FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR cmd/steganography-task

RUN go build

WORKDIR /app