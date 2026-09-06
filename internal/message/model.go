package message

type ResendEmailReceivedPayload struct {
	Type      string          `json:"type"`
	CreatedAt string          `json:"created_at"`
	Data      ResendEmailData `json:"data"`
}

type ResendEmailData struct {
	EmailId     string                   `json:"email_id"`
	CreatedAt   string                   `json:"created_at"`
	From        string                   `json:"from"`
	To          []string                 `json:"to"`
	Bcc         []string                 `json:"bcc"`
	Cc          []string                 `json:"cc"`
	ReceivedFor []string                 `json:"received_for"`
	MessageId   string                   `json:"message_id"`
	Subject     string                   `json:"subject"`
	Attachments []ResendEmailAttachments `json:"attachments"`
}

type ResendEmailAttachments struct {
	Id                 string `json:"id"`
	Filename           string `json:"filename"`
	ContentType        string `json:"content_type"`
	ContentDisposition string `json:"content_disposition"`
	ContentId          string `json:"content_id"`
}
