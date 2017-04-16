#!/bin/bash
# Creates .rpm/.deb packages for various platforms.

set -eu -o pipefail

if [ -z ${GOPATH+x} ]; then
  echo 'Please set $GOPATH.'
  exit 1
fi

# Directory of this script
SOURCE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Temporary directory
WORK_DIR=$(mktemp -d)
if [[ ! "$WORK_DIR" || ! -d "$WORK_DIR" ]]; then
  exit 1
fi
function cleanup {      
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

# Determine version number
VERSION="$(git describe --match "v*")"
VERSION="${VERSION#v}"

(
  cd "$SOURCE_DIR/static"
  webpack
)

build() {
  OS="$1"
  GOARCH="$2"
  PKGARCH="$3"
  TYPE="$4"

  ARCH_DIR="${WORK_DIR}/${OS}_${PKGARCH}_${TYPE}"
  mkdir "$ARCH_DIR"

  (
    cd "$ARCH_DIR"
    GOOS="$OS" GOARCH="$GOARCH" go build github.com/fxkr/weblabel/cmd/weblabel
  )

  fpm \
    --name "weblabel" \
    --description "A web interface for label printers" \
    --version "$VERSION" \
    --architecture "$PKGARCH" \
    --maintainer "Felix Kaiser <felix.kaiser@fxkr.net>" \
    --vendor "Felix Kaiser <felix.kaiser@fxkr.net>" \
    --url "https://github.com/fxkr/weblabel" \
    -t "$TYPE" \
    -s dir \
    --deb-systemd "$SOURCE_DIR/weblabel.service" \
    "$ARCH_DIR/weblabel=/usr/bin/weblabel" \
    "$SOURCE_DIR/config.yml.default=/usr/share/weblabel/config.yml" \
    "$SOURCE_DIR/weblabel.service=/usr/lib/systemd/system/weblabel.sevice" \
    "$SOURCE_DIR/static/dist/=/usr/share/weblabel/static"
}

build linux amd64 amd64 rpm
build linux amd64 amd64 deb
build linux arm armhf deb
