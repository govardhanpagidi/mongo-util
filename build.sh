#!/bin/sh
set -x
while getopts x:p:d:c:t:k:b:r:q: option
do
  case "${option}" in
      x) command=${OPTARG};;
      p) project_name=${OPTARG};;
      d) db=${OPTARG};;
      c) cluster=${OPTARG};;
      t) collection=${OPTARG};;
      k) data_api_key=${OPTARG};;
      b) atlas_pub_key=${OPTARG};;
      r) atlas_private_key=${OPTARG};;
      q) query=${OPTARG}
  esac
done

os=darwin
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

atlas_private_key="561002ad-65b1-4ac3-90fd-600a46218a39"
./mongo_util_${os}_${arch} -command="${command}" -project_name="${project_name}" \
  -db="${db}" -cluster="${cluster}" -collection="${collection}" -data_api_key="${data_api_key}" \
   -atlas_pub_key="${atlas_pub_key}" -atlas_private_key="${atlas_private_key}" -query=${query}

if [ $? -ne 0 ];
then
  echo "Job execution failed"
  exit 1
fi
echo "Job is executed"
