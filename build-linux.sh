#!/bin/sh

projectName=$1
set -x

os=linux
arch=amd64

cd ./src


# unzip go environment
go env
if [ $? -ne 0 ];
then
    echo "go not found, installing.."

    wget -c https://golang.org/dl/go1.17.linux-amd64.tar.gz
    tar -S -C /usr/local -xzf go1.17.linux-amd64.tar.gz
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

pwd

./mongo_util_linux_amd64 ${projectName}

echo "Job is executed"
