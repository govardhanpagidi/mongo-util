#!/bin/sh


projectName=$1

cd $(pwd)/src

$(pwd)/mongo_util_linux_amd64 $projectName

echo "execution completed."
