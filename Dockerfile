# 编译镜像
FROM golang:1.20-bullseye as build
ENV TZ=Asia/Shanghai
ENV DEBIAN_FRONTEND=noninteractive

# 构建依赖
RUN apt-get install apt-transport-https ca-certificates -y && \
    mv /etc/apt/sources.list /etc/apt/sources.list.bak && \
    echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free' >> /etc/apt/sources.list && \
    echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free' >> /etc/apt/sources.list && \
    echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free' >> /etc/apt/sources.list && \
    echo 'deb https://mirrors.tuna.tsinghua.edu.cn/debian-security bullseye-security main contrib non-free' >> /etc/apt/sources.list && \
    apt-get update && \
    apt-get install -yq git unzip tree wget tar && \
    apt-get clean && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false && \
    rm -rf /var/lib/apt/lists/*

# 获取文件
RUN mkdir -p /source
WORKDIR /source
COPY . .

# 编译
RUN bash build.sh

FROM alpine:latest

WORKDIR /data/apps/

# 收集数据
COPY --from=build /build ./feed-bundle

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    chmod +x ./feed-bundle/http && \
    chmod +x ./feed-bundle/grpc
