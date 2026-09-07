package message

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"

	"inbeux/internal/apperrors"
	userdb "inbeux/internal/user/db"

	"github.com/resend/resend-go/v4"
)

type UserProvider interface {
	CreateUser(context.Context, string) (userdb.User, error)
	GetUserByEmail(context.Context, string) (userdb.User, error)
}

func EmailReceiverHandler(messageService *MessageService, userService UserProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error reading body: ", "msg", err.Error())
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		headers := resend.WebhookHeaders{
			Id:        r.Header.Get("svix-id"),
			Timestamp: r.Header.Get("svix-timestamp"),
			Signature: r.Header.Get("svix-signature"),
		}

		err = client.Webhooks.Verify(&resend.VerifyWebhookOptions{
			Payload:       string(body),
			Headers:       headers,
			WebhookSecret: os.Getenv("RESEND_WEBHOOK_SECRET"),
		})

		if err != nil {
			slog.Error("webhook verification failed", "msg", err.Error())
			http.Error(w, "webhook verification failed", http.StatusBadRequest)
			return
		}

		var emailPayload ResendEmailReceivedPayload

		err = json.Unmarshal(body, &emailPayload)
		if err != nil {
			slog.Error("error parsing email payload json", "error", err)
			http.Error(w, "invalid json payload", http.StatusBadRequest)
			return
		}

		user, err := userService.CreateUser(r.Context(), emailPayload.Data.From)

		if err != nil {
			if errors.Is(err, apperrors.ErrUserAlreadyExists) {
				user, err = userService.GetUserByEmail(r.Context(), emailPayload.Data.From)
				if err != nil {
					slog.Error("failed to retrieve user", "msg", err.Error())
					http.Error(w, "failed to retrieve user", http.StatusInternalServerError)
					return
				}

			} else {
				slog.Error("failed to create user", "msg", err.Error())
				http.Error(w, "failed to create user", http.StatusInternalServerError)
				return
			}
		}

		err = messageService.ReceiveEmail(r.Context(), user.ID, emailPayload.Data.EmailId)

		if err != nil {
			slog.Error("failed to add email to database", "msg", err.Error())
			http.Error(w, "failed to add email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

}
