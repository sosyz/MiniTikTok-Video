#!/bin/bash

root=$(pwd)
cd /
PROTOC_ZIP=protoc-21.12-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v21.12/$PROTOC_ZIP
unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP

cd $root
cd Feed/proto
protoc *.proto --go_out=./ --go-grpc_out=./

cd $root/Feed/cmd/http
go build