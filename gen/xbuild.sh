#!/bin/sh
# Generates builds for all supported platforms.
# Based on the build script of https://github.com/gokcehan/lf

if [ -d ./dist ]; then
    rm -rf ./build || exit 5
fi

mkdir ./dist || exit 5

ERRORS=

build() {
    if [ $# != 2 ]; then
        echo "build function was called with inavlid parameters."
        ERRORS=1
        return 10
    fi

    GOOS="$1"
    GOARCH="$2"
    OUTFILE="dotm-$GOOS-$GOARCH"

    [ "$GOOS" = "windows" ] && OUTFILE="$OUTFILE.exe"

    printf "building for GOOS=%s and GOARCH=%s." "$GOOS" "$GOARCH"
    CGO_ENABLED=0 go build -o "dist/$OUTFILE"

    if [ "$?" != "0" ]; then
        ERRORS=1
        echo " errors."
    else
        echo " success."
    fi
    
    unset GOOS
    unset GOARCH
    unset OUTFILE
}

build darwin amd64
build darwin arm64
build linux amd64
build linux arm64
build windows amd64
build windows arm64

if [ -n "$ERRORS" ]; then
    printf "\nxbuild.sh: some targets faild to compile."
    exit 1
fi
