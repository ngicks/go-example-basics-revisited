# syntax=docker/dockerfile:1

ARG TAG_GOVER="1.25.0"
ARG TAG_DISTRO="bookworm"

FROM docker.io/library/golang:${TAG_GOVER}-${TAG_DISTRO} AS builder

ARG CGO_ENABLED="0"
ARG GOCACHE="/root/.cache/go-build"
ARG GOENV="/root/.config/go/env"
ARG GOPATH="/go"
ARG GOPRIVATE=""

ARG MAIN_PKG_PATH="."

ARG HTTP_PROXY
ARG HTTPS_PROXY=${HTTP_PROXY}
ARG NO_PROXY 
ARG http_proxy=${HTTP_PROXY}
ARG https_proxy=${HTTP_PROXY}
ARG no_proxy=${NO_PROXY}

# for curl, etc.
ARG SSL_CERT_FILE="/etc/ssl/certs/ca-certificates.crt"
ARG NODE_EXTRA_CA_CERTS=${SSL_CERT_FILE}
ARG DENO_CERT=${SSL_CERT_FILE}

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \ 
    --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt \
<<EOF
    rm -f /etc/apt/apt.conf.d/docker-clean
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
    apt-get update
    apt-get install -yqq --no-install-recommends git-lfs
EOF

WORKDIR /app/src

RUN --mount=type=secret,id=netrc,target=/root/.netrc \
    --mount=type=secret,id=goenv,target=/root/.config/go/env \
    --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt \
    --mount=type=cache,target=/go \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=bind,target=/app/src \
<<EOF
    go mod download
    # go generate ./...
    go build -o ../bin ${MAIN_PKG_PATH}
EOF

WORKDIR /app

# arm64
FROM gcr.io/distroless/static-debian12@sha256:ed92139a33080a51ac2e0607c781a67fb3facf2e6b3b04a2238703d8bcf39c40
# amd64
# FROM gcr.io/distroless/static-debian12@sha256:6ceafbc2a9c566d66448fb1d5381dede2b29200d1916e03f5238a1c437e7d9ea

COPY --from=builder /app/bin /app/bin

ENTRYPOINT [ "/app/bin" ]
