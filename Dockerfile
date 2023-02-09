FROM golang:1.20-bullseye as mini_tiktok_feed_http_builder

RUN apt-get install apt-transport-https ca-certificates -y \
    && mv /etc/apt/sources.list /etc/apt/sources.list.bak \
    && echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian-security bullseye-security main contrib non-free' >> /etc/apt/sources.list \
    && apt-get update \
    && apt-get install unzip git make gcc libprotobuf-dev protobuf-compiler -y \
    && git clone --recurse-submodules https://github.com/sosyz/mini_tiktok_feed.git /mini_tiktok_feed

WORKDIR /mini_tiktok_feed
RUN chmod +x ./build.sh && ./build.sh


FROM alpine:latest

WORKDIR /cloudreve
COPY --from=mini_tiktok_feed_http_builder /cloudreve_backend/Feed/cmd/http/http ./http

RUN apk update \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && chmod +x ./http \
EXPOSE 5212

ENTRYPOINT ["./http"]