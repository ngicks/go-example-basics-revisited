FROM docker.io/library/ubuntu:noble-20251013

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt \
<<EOF
    rm -f /etc/apt/apt.conf.d/docker-clean
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
    apt-get update
    apt-get install -yqq --no-install-recommends build-essential
EOF
