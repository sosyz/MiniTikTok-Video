# 安装protoc
wget https://ghproxy.com/https://github.com/protocolbuffers/protobuf/releases/download/v22.0/protoc-22.0-linux-x86_64.zip
unzip protoc-22.0-linux-x86_64.zip -d /usr/local
rm protoc-22.0-linux-x86_64.zip

# 设置 go 镜像
go env -w GO111MODULE=on
go env -w  GOPROXY=https://goproxy.io,direct

# 安装构建工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

cd Feed/proto
protoc *.proto --go_out=./ --go-grpc_out=./
cd ../..
mkdir -p $GOPATH/src/github.com/sosyz/mini_tiktok_feed
mv ./* $GOPATH/src/github.com/sosyz/mini_tiktok_feed
cd $GOPATH/src/github.com/sosyz/mini_tiktok_feed

go build -o /build/http github.com/sosyz/mini_tiktok_feed/Feed/cmd/http/
go build -o /build/grpc github.com/sosyz/mini_tiktok_feed/Feed/cmd/grpc/