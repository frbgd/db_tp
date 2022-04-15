FROM golang:latest AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/server/main.go

FROM ubuntu:18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt -y update && \
    apt install -y wget gnupg lsb-release ca-certificates && \
    apt -y update

WORKDIR /usr/src/app
COPY --from=build /app/main .

EXPOSE 5000

ENV DBHOST postgres
ENV DBPORT 5432
ENV DBNAME postgres
ENV DBUSER postgres
ENV DBPASS mysecretpassword

CMD ./main
