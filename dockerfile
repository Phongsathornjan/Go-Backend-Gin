FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

RUN mkdir -p /root/config

COPY --from=builder /app/main .

COPY .env .env

EXPOSE 8080

ENV GIN_MODE=release

CMD ["./main"]