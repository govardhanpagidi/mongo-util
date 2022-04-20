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
    echo "go not found, installing.."

    wget -c https://golang.org/dl/go1.17.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin

    go env
    if [ $? -ne 0 ];
    then
      echo "Problem in installing exiting.."
      exit 1
    fi
fi

env GOOS=${os} GOARCH=${arch} go build -o mongo_util_${os}_${arch}n

echo "Build is successful."
