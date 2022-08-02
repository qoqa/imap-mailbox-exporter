FROM golang:1.18.5-alpine as build

WORKDIR /app
COPY . /app

RUN go mod download
RUN go build -o /imap-mailbox-exporter

FROM alpine:3.16.1

WORKDIR /

RUN addgroup -S app && adduser -S -H -G app app
USER app

COPY --from=build /imap-mailbox-exporter /usr/local/bin/imap-mailbox-exporter

EXPOSE 9101

ENTRYPOINT [ "/usr/local/bin/imap-mailbox-exporter" ]