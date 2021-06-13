# syntax=docker/dockerfile:1
FROM golang AS builder
WORKDIR /src/app
RUN go get -d -v github.com/gofiber/fiber/v2
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker-snowflake .

FROM alpine
LABEL maintainer="Kawin Viriyaprasopsook <kawin.vir@zercle.tech>"

ARG	timezone=Asia/Bangkok

ENV	LANG en_US.UTF-8
ENV	LC_ALL en_US.UTF-8
ENV	TZ $timezone

# Add config repositories
RUN	echo 'https://download.nus.edu.sg/mirror/alpine/latest-stable/main' > /etc/apk/repositories \
    && echo 'https://download.nus.edu.sg/mirror/alpine/latest-stable/community' >> /etc/apk/repositories \
    && echo '@edge https://download.nus.edu.sg/mirror/alpine/edge/main' >> /etc/apk/repositories \
    && echo '@testing https://download.nus.edu.sg/mirror/alpine/edge/testing' >> /etc/apk/repositories \
    && mkdir /run/openrc \
    && touch /run/openrc/softlevel

# Update OS
RUN	apk update && apk upgrade \
    && apk add --no-cache \
    openrc \
    tzdata \
    bash \
    bash-completion \
    net-tools \
    bind-tools \
    mtr \
    ca-certificates

# Change locale
RUN echo $timezone > /etc/timezone \
    && cp /usr/share/zoneinfo/$timezone /etc/localtime

# Create app dir
COPY --from=builder /src/app/docker-snowflake /app/
COPY build/docker/ /app/
RUN ln -sf /app/docker-snowflake /usr/local/bin/app \
    && chmod +x /app/docker-entrypoint.sh
WORKDIR /app

ENTRYPOINT ["/app/docker-entrypoint.sh"]

CMD ["app", "prefork=false"]

EXPOSE 8080
