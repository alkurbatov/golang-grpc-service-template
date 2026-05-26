FROM golang:1.26.1 AS build

WORKDIR /build

ENV DEBIAN_FRONTEND=noninteractive

COPY Makefile go.mod go.sum ./
RUN make download

COPY . .
RUN make build

FROM ubuntu:26.04

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive \
    TZ=Europe/Moscow

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        locales && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    locale-gen en_US.UTF-8 && \
    update-locale LANG=en_US.UTF-8

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timez

COPY --from=build /build/dist .

RUN chown -R nobody:nogroup /app \
    && chmod -R 700 /app
USER nobody

ENTRYPOINT ["./templatesrv"]
