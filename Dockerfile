FROM golang:1.20-alpine as build
WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /ytdl ./cmd/web

FROM python:3.11-alpine
WORKDIR /app
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

COPY --from=build /ytdl ./ytdl
COPY ui ./ui

EXPOSE 8080
ENTRYPOINT ["./ytdl", "-addr=:8080", "-dir=/dl" ]
