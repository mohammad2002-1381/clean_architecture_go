FROM golang:1.26-alpine AS builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN swag init -g ./cmd/api/main.go -d . --parseDependency --parseInternal -o ./cmd/docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./cmd/api/main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/api-server .
COPY --from=builder /app/.env .env

EXPOSE 8000
CMD ["./api-server"]
