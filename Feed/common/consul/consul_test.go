package consul

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
)

func TestRegister(t *testing.T) {
	cfg := &conf.ContainerConfig{
		AppName:            "feed",
		AppVersion:         "v1",
		ENV:                "dev",
		RegionName:         "local",
		ListenHost:         "127.0.0.1",
		ListenPort:         8080,
		ServiceDiscoverUrl: "http://127.0.0.1:8500",
		AppRoot:            "/Users/sosyz/go/src/github.com/sosyz/mini_tiktok_feed/Feed",
		CFGAccessKey:       "xxxx",
		ImageTag:           "feed:v1",
	}
	go func() {
		http.ListenAndServe(":8080", nil)
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "ok")
		})
	}()

	t.Logf("%v\n", cfg)
	consulConn, err := ConnectConsul(cfg)
	if err != nil {
		t.Error(err)
	}
	// defer consulConn.Close()
	svrID, err := RegisterService(cfg, consulConn)
	if err != nil {
		t.Error(err)
	}
	t.Log(svrID)
	err = UnregisterService(svrID, consulConn)
	if err != nil {
		t.Error(err)
	}
}
