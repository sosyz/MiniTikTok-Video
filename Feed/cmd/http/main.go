package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
)

func GetServiceConnectInfo(consulUrl string, needServices ...any) map[string]*conf.Server {
	ret := make(map[string]*conf.Server)
	for idx, needneedService := range needServices {
		ret[needneedService.(string)] = &conf.Server{
			Host: "127.0.0.1",
			Port: 8080 + idx,
		}
	}
	return ret
}

func main() {
	consulUrl := os.Getenv("DREAM_SERVICE_DISCOVERY_URI")
	serviceMap := GetServiceConnectInfo(consulUrl, "feed", "auth")
	InitService(serviceMap["feed"], serviceMap["auth"])

	r := InitRouter("")
	service := ""
	fmt.Print(service)
	r.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	r.Run(":8080")
}
