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

# let gpg key unlocked for ssh login.
# As you can see in https://github.com/containers/buildah/blob/v1.42.1/pkg/sshagent/sshagent.go#L126-L131
# buildah sets 2 sec timeout for ssh-agent so you have low chance to successfully enter passphrase.
ssh-add -T ~/.ssh/id_ecdsa.pub

# this is really needed.
export HTTP_PROXY=${HTTP_PROXY}
export HTTPS_PROXY=${HTTPS_PROXY:-$HTTP_PROXY}
# maybe being empty is ok.
export NO_PROXY=${NO_PROXY:-""}
export http_proxy=${http_proxy:-$HTTP_PROXY}
export https_proxy=${https_proxy:-$HTTPS_PROXY}
export no_proxy=${no_proxy:-$NO_PROXY}

podman buildx build \
    --platform linux/${arch} \
    --build-arg TAG_GOVER=${TAG_GOVER} \
    --build-arg MAIN_PKG_PATH=${MAIN_PKG_PATH:-./} \
    --build-arg GOPRIVATE=${GOPRIVATE:-""} \
    --secret id=netrc,src=${NETRC:-$HOME/.netrc} \
    --secret id=goenv,src=$(go env GOENV) \
    --build-arg SSL_CERT_FILE=${SSL_CERT_FILE:-/etc/ssl/certs/ca-certificates.crt} \
    --secret id=certs,src=${SSL_CERT_FILE:-/etc/ssl/certs/ca-certificates.crt} \
    --secret id=HTTP_PROXY \
    --secret id=HTTPS_PROXY \
    --secret id=NO_PROXY \
    --secret id=http_proxy \
    --secret id=https_proxy \
    --secret id=no_proxy \
    -t ${1}-${arch} \
    -f Containerfile \
    .
