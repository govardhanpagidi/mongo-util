#!/bin/sh

#archs=(amd64 arm64 ppc64le ppc64 s390x)
#archs=(amd64)
#os=linux

set -x
arch=amd64
os=darwin
cd ./src

go env
if [ $? -ne 0 ];
then
    echo "setting go path."
#    wget -c https://golang.org/dl/go1.17.linux-amd64.tar.gz
#    tar -S -C /usr/local -xzf go1.17.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    go env
    if [ $? -ne 0 ];
    then
      echo "Problem in installing exiting.."
      exit 1
    fi
fi

env GOOS=${os} GOARCH=${arch} go build -o mongo_util_${os}_${arch}

echo "Build is successful."
