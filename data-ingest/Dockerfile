FROM golang:1.21 AS builder
LABEL authors="ginos"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app cmd/main.go

CMD ["./main"]