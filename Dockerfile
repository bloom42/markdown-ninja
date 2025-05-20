####################################################################################################
## Build webapp & themes
####################################################################################################
# node:18
FROM --platform=$BUILDPLATFORM node@sha256:867be01f97d45cb7d89a8ef0b328d23e8207412ebec4564441ed8cabc8cc4ecd AS builder_js

RUN apt update && apt upgrade -y && \
    apt install -y libimage-exiftool-perl make

WORKDIR /mdninja
COPY . ./

# build webapp
WORKDIR /mdninja/webapp
RUN make exif
RUN make clean
RUN make install_ci
RUN make build

# build themes
WORKDIR /mdninja/themes/blog
RUN make exif
RUN make clean
RUN make install_ci
RUN make build

WORKDIR /mdninja/themes/docs
RUN make exif
RUN make clean
RUN make install_ci
RUN make build

####################################################################################################
## Build mdninja and mdninja-server
####################################################################################################
# golang:1.24
FROM golang@sha256:39d9e7d9c5d9c9e4baf0d8fff579f06d5032c0f4425cdec9e86732e8e4e374dc AS go

FROM ubuntu:24.04 AS builder_go


ENV TZ="UTC"
# ENV LC_ALL="en_US.UTF-8"
# ENV LANG="en_US.UTF-8"
# ENV LANGUAGE="en_US:en"

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apt update && apt upgrade -y && \
    apt install -y ca-certificates git make tzdata binutils mailcap libcap2-bin
RUN update-ca-certificates

# setup go
COPY --from=go /usr/local/go /usr/local/go
ENV GOROOT="/usr/local/go"
ENV PATH="$PATH:$GOROOT/bin"
ENV GOPROXY=direct
ENV GOTOOLCHAIN="local"
RUN go telemetry off


WORKDIR /mdninja
COPY . ./

RUN make clean

COPY --from=builder_js /mdninja/webapp/dist/ /mdninja/webapp/dist/
COPY --from=builder_js /mdninja/themes/blog/dist/ /mdninja/themes/blog/dist/
COPY --from=builder_js /mdninja/themes/docs/dist/ /mdninja/themes/docs/dist/

# download_deps is disabled because it drastically slowed down builds
# RUN make download_deps
RUN make mdninja-server
RUN make mdninja
# RUN make verify_deps

# This allows us to listen on ports 1024, to bind a TLS listener on port 443 for example.
# RUN setcap CAP_NET_BIND_SERVICE=+eip /mdninja/dist/mdninja-server

####################################################################################################
## This container is used to get the correct files to scratch
####################################################################################################
FROM --platform=$BUILDPLATFORM debian:12-slim AS builder_files

# appuser
ENV USER=mdninja
ENV UID=10001

# mailcap is used for content type (MIME type) detection
# tzdata is used for timezones info
RUN apt update && apt upgrade -y && \
    apt install -y mailcap ca-certificates adduser wget tzdata
RUN update-ca-certificates

ENV TZ="UTC"
RUN echo "${TZ}" > /etc/timezone

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

####################################################################################################
## Final image
####################################################################################################

# See https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
# for more information about how to create Go containers FROM scratch
FROM scratch

ENV TZ="UTC"
ENV LC_ALL="en_US.UTF-8"
ENV LANG="en_US.UTF-8"
ENV LANGUAGE="en_US:en"

# /etc/nsswitch.conf may be used by some DNS resolvers
# /etc/mime.types may be used to detect the MIME type of files
COPY --from=builder_files --chmod=444 \
    /etc/passwd \
    /etc/group \
    /etc/nsswitch.conf \
    /etc/mime.types \
    /etc/timezone \
    /etc/


COPY --from=builder_files --chmod=444 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder_files --chmod=444 /usr/share/zoneinfo /usr/share/zoneinfo

# Copy our builds
COPY --from=builder_go /mdninja/dist/mdninja /usr/local/bin/mdninja
COPY --from=builder_go /mdninja/dist/mdninja-server /usr/local/bin/mdninja-server

# Use an unprivileged user.
USER mdninja:mdninja

# the scratch image doesn't have a /tmp folder so we need to create it
WORKDIR /tmp

# Final working directory
WORKDIR /mdninja

ENTRYPOINT ["/usr/local/bin/mdninja"]

EXPOSE 8080
