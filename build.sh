#!/bin/sh

#archs=(amd64 arm64 ppc64le ppc64 s390x)
#archs=(amd64)
#os=linux

set -x
arch=amd64
os=darwin
cd ./src


# unzip go environment
go env
if [ $? -ne 0 ];
then
    echo "go not found, installation started..."
    go_env="go1.17.darwin-amd64.tar.gz"
    curl -L -o go.pkg https://go.dev/dl/go1.17.${os}-${arch}.pkg
    echo "go package downloaded."
    #rm -rfv /usr/local/go
    open -S go.pkg

    #tar -zxf $go_env go

    #open -S go.pkg
    go env
    if [ $? -ne 0 ];
    then
        echo "fail in extract go"
        exit 1
    fi
    echo "OK for extract go"
    rm -rf $go_env

    # prepare PATH, GOROOT and GOPATH
    #export PATH=$(pwd)/go/bin:$PATH
    export PATH=$PATH:$GOPATH/bin
    export GOROOT=$(pwd)/go
    export GOPATH=$(pwd)
fi
pwd

env GOOS=${os} GOARCH=${arch} go build -o mongo_util_${os}_${arch}

echo "Build is successful."
