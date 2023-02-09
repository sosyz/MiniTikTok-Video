FROM golang:1.20-bullseye as mini_tiktok_feed_http_builder

RUN git clone --recurse-submodules https://github.com/sosyz/mini_tiktok_feed.git /mini_tiktok_feed

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