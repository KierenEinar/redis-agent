#! /bin/bash

pname=`whoami`

if [ -z $GOPATH ]
then
    echo "system not set GOPATH, export GOPATH to ~$pname/go"
    export GOPATH=~$pname/go
fi

cd $GOPATH/src/redis-agent

echo "===================================start build========================================="

go clean && glide install && go build

echo "===================================build successful========================================="
