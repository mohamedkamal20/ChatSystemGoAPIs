# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /chatSystemGoAPIs

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8085

CMD ["./main"]