FROM golang:1.13-stretch AS builder

# Building project
WORKDIR /build

COPY . .
RUN go build -v -o main ./main.go

FROM ubuntu:19.04
ENV DEBIAN_FRONTEND=noninteractive

# Expose server & database ports
EXPOSE 5000
EXPOSE 5432

RUN apt-get update && apt-get install -y postgresql-11

USER postgres

# Create & configure database
COPY ./assets/db/db.sql .
RUN service postgresql start &&\
	psql --command "CREATE USER forum_user WITH SUPERUSER PASSWORD 'forum_password';" &&\
	createdb -O forum_user forum_db &&\
    psql -f ./db.sql -d forum_db &&\
    service postgresql stop
ENV DATABASE_URL=postgres://forum_user:forum_password@localhost/forum_db

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /app

# Copying built binary
COPY --from=builder /build/main .
CMD service postgresql start && ./main