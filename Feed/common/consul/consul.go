package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
)

func ConnectConsul(cfg *conf.EnvConfig) (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	config.Address = cfg.ServiceDiscoverUrl
	return consulapi.NewClient(config)
}

func RegisterService(cfg *conf.EnvConfig, consulServer *consulapi.Client) error {
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = cfg.ImageTag
	registration.Name = cfg.AppName
	registration.Port = cfg.ListenPort
	registration.Tags = []string{cfg.AppVersion, cfg.ENV, cfg.RegionName}
	registration.Address = cfg.ListenHost
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", cfg.ListenHost, cfg.ListenPort),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	return consulServer.Agent().ServiceRegister(registration)
}

func GetServiceConnectInfo(consulServer *consulapi.Client, needServices ...any) map[string]*conf.Server {
	ret := make(map[string]*conf.Server)
	for _, needService := range needServices {
		service, _, err := consulServer.Health().Service(needService.(string), "", true, nil)
		if err != nil {
			panic(err)
		}
		if len(service) == 0 {
			panic("no service found")
		}

		ret[needService.(string)] = &conf.Server{
			Host: service[0].Service.Address,
			Port: service[0].Service.Port,
		}
	}
	return ret
}
