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
    echo "go not found, exiting"
    exit 1
fi
pwd

env GOOS=${os} GOARCH=${arch} go build -o mongo_util_${os}_${arch}

echo "Build is successful."
