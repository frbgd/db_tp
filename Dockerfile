# Step 1. build_step
FROM golang:1.14-stretch AS build_step
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o app api_service/main.go

# Step 2. release_step
FROM ubuntu:18.04 AS release_step

ENV DEBIAN_FRONTEND=noninteractive
ENV PGVER 12

RUN apt -y update && \
    apt install -y wget gnupg && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" >> /etc/apt/sources.list.d/pgdg.list && \
    apt -y update

RUN apt -y update && apt install -y \
        postgresql-$PGVER \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build_step /app/db/init.sql /db/init.sql

USER postgres
RUN service postgresql start && \
    psql --command "ALTER USER postgres PASSWORD 'mysecretpassword';" && \
    psql postgres --echo-all --file /db/init.sql && \
    service postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf && \
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

COPY --from=build_step /app/api_service /app/

WORKDIR /app

EXPOSE 5432
EXPOSE 5000

CMD service postgresql start && ./app
