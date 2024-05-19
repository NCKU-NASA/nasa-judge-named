FROM golang:1.20 as builder

RUN apt-get update && apt-get full-upgrade -y && apt-get install make -y

WORKDIR /src
COPY go.mod .
COPY go.sum .
ARG CACHE
RUN --mount=type=cache,target="$CACHE" go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get
COPY . .
RUN --mount=type=cache,target="$CACHE" make clean && make

FROM debian:latest as release

RUN apt-get update && apt-get full-upgrade -y && apt-get install ca-certificates curl wget wireguard jq bind9 bind9utils dnsutils procps iputils-ping mtr -y && update-ca-certificates

COPY --from=builder /src/bin /app
COPY docker-entrypoint.sh /app
COPY default.named /etc/default/named

RUN chmod +x /app/docker-entrypoint.sh && \
    mkdir -p /zone && \
    mkdir -p /var/cache/bind/ && \
    chown bind:bind /var/cache/bind/ && \
    mkdir -p /var/log/bind && \
    chown bind:bind /var/log/bind

WORKDIR /app

RUN touch .env
RUN touch init.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
