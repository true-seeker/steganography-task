FROM golang:latest

WORKDIR /app

#COPY go.mod .
#COPY go.sum .
COPY . .

RUN go mod download

WORKDIR cmd/it_planet_task

RUN go build
