FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /acore ./main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates wget
COPY --from=builder /acore /acore
COPY --from=builder /app/views /views
COPY --from=builder /app/.env /.env
EXPOSE 8080
ENTRYPOINT ["/acore"]
