FROM ubuntu:18.04

ENV DEBIAN_FRONTEND=noninteractive
ENV PGVER 12

RUN apt -y update && \
    apt install -y wget gnupg && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" >> /etc/apt/sources.list.d/pgdg.list && \
    apt -y update

RUN apt -y update && apt install -y \
        postgresql-$PGVER \
    && rm -rf /var/lib/apt/lists/* \

RUN cd /tmp && wget https://golang.org/dl/go1.17.5.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz

COPY ./db /db

USER postgres
RUN service postgresql start && \
    psql --command "ALTER USER postgres PASSWORD 'mysecretpassword';" && \
    psql postgres --echo-all --file /db/init.sql && \
    service postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf && \
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

WORKDIR /usr/src/app
COPY api_service .

EXPOSE 5432
EXPOSE 5000

ENV DBHOST 127.0.0.1
ENV DBPORT 5432
ENV DBNAME postgres
ENV DBUSER postgres
ENV PGPASSWORD mysecretpassword

CMD service postgresql start && /usr/local/go/bin/go build main.go && ./main
