FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /bin/app ./
COPY --from=builder /app/.env .env
COPY --from=builder /app/public ./public
COPY --from=builder /app/styles ./styles

EXPOSE 5555

CMD ["./app"]
