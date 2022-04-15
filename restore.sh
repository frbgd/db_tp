export PGPASSWORD=mysecretpassword
psql -U postgres
postgres=# create database newdb;
postgres=# \c newdb
DROP DATABASE postgres;
CREATE DATABASE postgres;
\q
pg_restore -h localhost -p 5432 -U postgres -d postgres /root/bazka.sql