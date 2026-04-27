FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-s -w" github.com/pressly/goose/v3/cmd/goose@v3.22.1

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata wget

COPY --from=builder /server /usr/local/bin/server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations /migrations

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["server"]
