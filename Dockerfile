FROM ubuntu:18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y software-properties-common gcc vim && \
    add-apt-repository -y ppa:deadsnakes/ppa
RUN apt-get update && apt-get install -y python3.9 python3.9-dev python3-distutils python3-pip python3-apt build-essential

ENV PGVER 12

RUN apt -y update && \
    apt install -y wget gnupg && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" >> /etc/apt/sources.list.d/pgdg.list && \
    apt -y update

RUN apt -y update && apt install -y \
        postgresql-$PGVER \
    && rm -rf /var/lib/apt/lists/*

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

COPY ./api_service/requirements.txt /app/requirements.txt

RUN python3 -m pip install --upgrade pip && pip3 install --upgrade setuptools && pip3 install -r /app/requirements.txt

COPY ./api_service /app

WORKDIR /app

EXPOSE 5432
EXPOSE 5000

CMD service postgresql start && gunicorn main:app --workers $(($(nproc) + 1)) --worker-class uvicorn.workers.UvicornWorker --bind 0.0.0.0:5000
