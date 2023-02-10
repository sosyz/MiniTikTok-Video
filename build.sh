root=$(pwd)
file=protoc-21.12-linux-x86_64.zip

curl -OL https://ghproxy.com/https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip
unzip protoc-21.12-linux-x86_64.zip -d /usr/local
rm protoc-21.12-linux-x86_64.zip

tree ./
cd $root/mini_tiktok_feed/Feed/proto/
protoc *.proto --go_out=./ --go-grpc_out=./ \
go build -o $root/mini_tiktok_feed/Feed/cmd/http/http