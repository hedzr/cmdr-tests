#!/usr/bin/env bash

is_darwin() { [[ $OSTYPE == *darwin* ]]; }

is_darwin && LS_OPTS="-G" || LS_OPTS="--color"

test_docgo() {
  local DOCGO="${1:-doc.go}"
  if [ -f "$DOCGO" ]; then
    if grep -E "Version[ \t]+=[ \t]+" $DOCGO; then
      :
    else
      false
    fi
  else
    false
  fi
}

build-app() {
  local name="${1:-blueprint}"
  local SRCS="./examples/${name}/"
  local AN="${name}"

  W_PKG="${W_PKG:-github.com/hedzr/cmdr/v2/conf}"
  GOPROXY="${GOPROXY:-https://goproxy.io,direct}"
  DOCGO="doc.go"
  test_docgo "$DOCGO" || {
    DOCGO="./_examples/doc.go" && test_docgo "$DOCGO" || {
      DOCGO="./examples/doc.go" && test_docgo "$DOCGO" || {
        DOCGO="./examples/demo/doc.go" && test_docgo "$DOCGO" || {
          echo "CANNOT locate where is the Version."
          exit -1
        }
      }
    }
  }
  V1="$(grep -E "Version[ \t]+=[ \t]+" doc.go|grep -Eo "[0-9.]+")"
  VERSION="${VERSION:-$V1}"
  GIT_REVISION="$(git rev-parse --short HEAD)"
  GIT_SUMMARY="$(git describe --tags --dirty --always)"
  GIT_DESC="$(git log --oneline -1)"
  GOVERSION="$(go version)"
  BUILDTIME="$(date -u '+%Y-%m-%d_%H-%M-%S')"
  BUILDTIMESTAMP="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  LDFLAGS="-s -w \
    -X '${W_PKG}.Githash=$GIT_REVISION' \
    -X '${W_PKG}.GitSummary=$GIT_SUMMARY' \
    -X '${W_PKG}.GitDesc=$GIT_DESC' \
    -X '${W_PKG}.BuilderComments=$BUILDER_COMMENT' \
    -X '${W_PKG}.GoVersion=$GOVERSION' \
    -X '${W_PKG}.Version=$VERSION' \
    -X '${W_PKG}.AppName=$APPNAME'
    -X '${W_PKG}.Buildstamp=$BUILDTIMESTAMP'"

  echo;echo;echo "Using GOPROXY: $GOPROXY";echo "    CN: $CN";echo

  GOOS="${GOOS:-linux}"
  CGO_ENABLED="${CGO_ENABLED:-0}"

  go build -v -a -installsuffix cgo \
    -ldflags "$LDFLAGS" \
    -o bin/$AN $SRCS && \
  chmod +x bin/$AN && \
  ls -la $LS_OPTS bin/$AN && \
  echo "version" && ./bin/$AN --version && \
  echo "build-info" && ./bin/$AN --built-info
  # echo
  # echo "$LDFLAGS"
}

cmd="${1:-build-app}"
(($#)) && shift
$cmd "$@"