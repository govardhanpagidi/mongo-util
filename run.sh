#!/bin/sh


projectName=$1

cd $(pwd)/src

$(pwd)/mongo_util_darwin_amd64 $projectName

echo "execution completed."
