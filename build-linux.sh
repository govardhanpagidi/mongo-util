#!/bin/sh

projectName=$1
set -x

os=linux
arch=amd64

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

go mod vendor
env GOOS=${os} GOARCH=${arch} go build -o mongo_util_${os}_${arch}
if [ $? -ne 0 ];
then
  echo "Build is failed."
  exit 1
fi
echo "Build is successful."

pwd

./mongo_util_linux_amd64 ${projectName}

if [ $? -ne 0 ];
then
  echo "Job execution failed"
  exit 1
fi
echo "Job is executed"
