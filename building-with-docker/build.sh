#! /bin/sh

podman buildx build \
    --build-arg HTTP_PROXY=${HTTP_PROXY}\
    --build-arg HTTPS_PROXY=${HTTPS_PROXY}\
    --build-arg MAIN_PKG_PATH=${MAIN_PKG_PATH:-./}\
    --build-arg GOPRIVATE=${GOPRIVATE}\
    --secret id=certs,src=/etc/ssl/certs/ca-certificates.crt\
    --secret id=goenv,src=$(go env GOENV)\
    --ssh default=${SSH_AUTH_SOCK} \
    -t $1\
    -f Dockerfile\
    .
