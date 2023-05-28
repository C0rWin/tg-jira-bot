FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY .env .
ENTRYPOINT ["/app/app"]
