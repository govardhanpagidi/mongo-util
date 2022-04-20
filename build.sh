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
   yum install golang -y
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
