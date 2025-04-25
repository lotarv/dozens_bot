FROM golang:1.24-alpine3.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk add --no-cache bash

EXPOSE 3001