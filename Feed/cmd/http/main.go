package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	cfg "github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	cs "github.com/sosyz/mini_tiktok_feed/Feed/common/consul"
)

func main() {
	conf := cfg.ReadContainerConfig()
	fmt.Printf("%v\n", conf)
	consulConn, err := cs.ConnectConsul(conf)
	if err != nil {
		panic(err)
	}
	dependService := cs.GetServiceConnectInfo(consulConn, "bawling-minitiktok-feed", "bawling-minitiktok-auth")

	cs.RegisterService(conf, consulConn)
	s3Conf := cfg.ReadS3ConfigByEnv()
	volConf := cfg.ReadSecretByEnv("VOL")
	InitService(dependService["feed"], dependService["auth"], s3Conf, volConf)

	r := InitRouter("")
	service := ""
	fmt.Print(service)
	r.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	cs.RegisterService(conf, consulConn)
	r.Run(fmt.Sprintf("%s:%d", conf.ListenHost, conf.ListenPort))
}
