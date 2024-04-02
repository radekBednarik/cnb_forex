# syntax=docker/dockerfile:1

FROM alpine:3.19.1 as base

RUN apk --no-cache add postgresql postgresql-contrib postgresql-client 

ARG DB_NAME
ARG DB_USER
ARG DB_PASSWORD

ENV DB_NAME=${DB_NAME}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  appuser
USER ${DB_NAME} 

RUN /etc/init.d/postgresql setup

CMD ["postgres", "-D", "/var/lib/postgresql/data", "-c", "config_file=/etc/postgresql/postgresql.conf]

