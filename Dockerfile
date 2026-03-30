FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/Project/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /bin/app ./app
COPY --from=builder /app/frontend ./frontend
COPY config.yaml .

EXPOSE 8080

CMD ["./app"]