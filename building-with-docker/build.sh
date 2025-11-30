#! /bin/sh

set -Cue

TAG_GOVER=1.25.0
if [ -f ./ver ]; then
  TAG_GOVER=$(cat ./ver)
fi

arch=${TARGET_ARCH:-$(go env GOARCH)}

echo $arch

# let gpg key unlocked for ssh login.
# As you can see in https://github.com/containers/buildah/blob/v1.42.1/pkg/sshagent/sshagent.go#L126-L131
# buildah sets 2 sec timeout for ssh-agent so you have low chance to successfully enter passphrase.
ssh-add -T ~/.ssh/id_ecdsa.pub

podman buildx build \
    --platform linux/${arch} \
    --build-arg TAG_GOVER=${TAG_GOVER} \
    --build-arg HTTP_PROXY=${HTTP_PROXY:-""} \
    --build-arg HTTPS_PROXY=${HTTPS_PROXY:-""} \
    --build-arg MAIN_PKG_PATH=${MAIN_PKG_PATH:-./} \
    --build-arg GOPRIVATE=${GOPRIVATE:-""} \
    --secret id=certs,src=/etc/ssl/certs/ca-certificates.crt \
    --secret id=goenv,src=$(go env GOENV) \
    --ssh default=${SSH_AUTH_SOCK:-""} \
    -t ${1}-${arch} \
    -f Containerfile \
    .
