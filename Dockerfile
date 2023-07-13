FROM golang:1.20-alpine as build
WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /ytdl ./cmd/web

FROM python:3.11-alpine
# install ffmpeg
RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

# install yt-dlp
WORKDIR /app
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

# copy binary and html template files
COPY --from=build /ytdl ./ytdl
COPY docker-entrypoint.sh ./docker-entrypoint.sh
COPY ui ./ui

EXPOSE 8080
ENTRYPOINT ["./docker-entrypoint.sh", "-addr=:8080", "-dir=/dl", "-zip=/dl/zip" ]
