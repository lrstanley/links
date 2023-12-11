# syntax = docker/dockerfile:1.4

# build image
FROM golang:latest as build
WORKDIR /build
COPY go.sum go.mod Makefile /build/
RUN make go-fetch
COPY . /build/
RUN make

# runtime image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=build /build/links /usr/local/bin/links

# runtime params
VOLUME /data
EXPOSE 80
WORKDIR /
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["/usr/local/bin/links", "--http", "0.0.0.0:80", "--behind-proxy", "--db", "/data/store.db"]
