
# build image
FROM golang:alpine as build
WORKDIR /build
RUN apk add --no-cache g++ make
COPY go.sum go.mod Makefile /build/
RUN make fetch
COPY . /build/
RUN make

# runtime image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/docker-library/golang/blob/1eb096131592bcbc90aa3b97471811c798a93573/1.14/alpine3.12/Dockerfile#L9
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=build /build/links /usr/local/bin/links

# runtime params
VOLUME /data
EXPOSE 80
WORKDIR /
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["links", "--http", "0.0.0.0:80", "--behind-proxy", "--db", "/data/store.db"]