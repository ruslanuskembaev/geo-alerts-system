FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/geo-alerts ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/geo-alerts /app/geo-alerts

EXPOSE 8080
ENTRYPOINT ["/app/geo-alerts"]
