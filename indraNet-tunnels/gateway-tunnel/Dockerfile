FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o gateway-tunnel .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/gateway-tunnel .

EXPOSE 8082
EXPOSE 8080

CMD ["./gateway-tunnel"]