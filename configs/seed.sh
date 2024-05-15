#/bin/sh

psql -U postgres -d postgres -c "\COPY cpu_usage FROM docker-entrypoint-initdb.d/cpu_usage.csv CSV HEADER"