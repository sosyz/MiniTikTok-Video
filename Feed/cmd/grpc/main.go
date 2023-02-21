package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	cs "github.com/sosyz/mini_tiktok_feed/Feed/common/consul"
	pb_feed "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/feed"
	pb_usr "github.com/sosyz/mini_tiktok_feed/Feed/proto/pb/user"
	srv_feed "github.com/sosyz/mini_tiktok_feed/Feed/service/feed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	go func() {
		http.ListenAndServe(":8080", nil)
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "ok")
		})
	}()

	cfg := conf.ReadContainerConfig()
	consulConn, err := cs.ConnectConsul(cfg)
	if err != nil {
		panic(err)
	}
	dependService := cs.GetServiceConnectInfo(consulConn, "user")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", dependService["user"].Host, dependService["user"].Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	userService := pb_usr.NewUserServiceClient(conn)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.ListenPort))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	feedService, err := srv_feed.NewFeedService(
		conf.ReadNeo4jConfigByEnv(),
		conf.ReadSecretByEnv("VOL"),
		userService,
		conf.ReadNodeConfigByEnv().Id,
	)
	if err != nil {
		panic(err)
	}
	pb_feed.RegisterFeedServiceServer(s, feedService)
	cs.RegisterService(cfg, consulConn)
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
