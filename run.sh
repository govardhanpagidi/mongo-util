#!/bin/sh

#archs=(amd64 arm64 ppc64le ppc64 s390x)
archs=(amd64)

# go build
for arch in ${archs[@]}
do
        env GOOS=linux GOARCH=${arch} go build -o util_${arch}
done

# Copying files to target machine
files=(util_${arch} config.json gcp.json )

for file in ${files[@]}
do
    scp -i "~/Downloads/zebra.pem" ${file} ec2-user@ec2-54-196-226-146.compute-1.amazonaws.com:$home
    #aws --region us-east-1 s3 cp ${file} s3://ops-test-builds/go/
done
echo "done"
