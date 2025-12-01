#! /bin/sh

set -Cue

if [ -z ${1:-""} ]; then
  echo "set repo:tag as first cli argument"
  exit 1
fi

TAG_GOVER=1.25.0
if [ -f ./ver ]; then
  TAG_GOVER=$(cat ./ver)
fi

arch=${TARGET_ARCH:-""}

if [ -z ${arch} ]; then
  case $(uname -m) in
    "x86_64")
      arch="amd64";;
    "x86_64-AT386")
      arch="amd64";;
    "aarch64_be")
      arch="arm64be";;
    "aarch64")
      arch="arm64";;
    "armv8b")
      arch="arm64";;
    "armv8l")
      arch="arm64";;
  esac
fi

if [ -z $arch ]; then
  echo "arch unknown: $(uname -m)"
  exit 1
fi

echo $arch

# this is really needed.
export HTTP_PROXY=${HTTP_PROXY}
export HTTPS_PROXY=${HTTPS_PROXY:-$HTTP_PROXY}
# maybe being empty is ok.
export NO_PROXY=${NO_PROXY:-""}
export http_proxy=${http_proxy:-$HTTP_PROXY}
export https_proxy=${https_proxy:-$HTTPS_PROXY}
export no_proxy=${no_proxy:-$NO_PROXY}

docker buildx build \
    --platform linux/${arch} \
    --build-arg TAG_GOVER=${TAG_GOVER} \
    --build-arg MAIN_PKG_PATH=${MAIN_PKG_PATH:-./} \
    --build-arg GOPRIVATE=${GOPRIVATE:-""} \
    --secret id=goenv,src=$(go env GOENV) \
    --secret id=netrc,src=${NETRC:-$HOME/.netrc} \
    --secret id=certs,src=${SSL_CERT_FILE:-/etc/ssl/certs/ca-certificates.crt} \
    --secret id=HTTP_PROXY,type=env \
    --secret id=HTTPS_PROXY,type=env \
    --secret id=NO_PROXY,type=env \
    --secret id=http_proxy,type=env \
    --secret id=https_proxy,type=env \
    --secret id=no_proxy,type=env \
    -t ${1}-${arch} \
    -f behind-proxy.Dockerfile \
    .
