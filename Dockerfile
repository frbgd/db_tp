FROM ubuntu:18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt -y update && \
    apt install -y wget gnupg lsb-release ca-certificates && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" >> /etc/apt/sources.list.d/pgdg.list && \
    apt -y update

RUN apt -y update && apt install -y \
        postgresql-14 \
    && rm -rf /var/lib/apt/lists/* \

RUN cd /tmp && wget https://golang.org/dl/go1.17.5.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz

COPY ./db /db

USER postgres
RUN service postgresql start && \
    psql --command "ALTER USER postgres PASSWORD 'mysecretpassword';" && \
    psql postgres --echo-all --file /db/init.sql && \
    service postgresql stop

RUN echo "local   all             postgres                                peer\nlocal   all             all                                     md5\nhost    all             all             127.0.0.1/32            scram-sha-256\nhost    all             all             0.0.0.0/0               md5" > /etc/postgresql/14/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/14/main/postgresql.conf
RUN echo "synchronous_commit = off\nfsync = off\nshared_buffers = 1GB\nwork_mem=512MB\nunix_socket_directories = '/var/run/postgresql'\nunix_socket_permissions = 0777" >> /etc/postgresql/14/main/postgresql.conf

USER root
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

WORKDIR /usr/src/app
COPY . .

EXPOSE 5432
EXPOSE 5000

ENV DBHOST 127.0.0.1
ENV DBPORT 5432
ENV DBNAME postgres
ENV DBUSER postgres
ENV PGPASSWORD mysecretpassword

CMD service postgresql start && /usr/local/go/bin/go build -o ./main main.go && ./main
