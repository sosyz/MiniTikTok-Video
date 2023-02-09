#!/bin/bash

root=$(pwd)
cd Feed/proto
protoc *.proto --go_out=./ --go-grpc_out=./

cd $root/Feed/cmd/http
go build