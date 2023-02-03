package visual

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/conf"
	"github.com/spf13/viper"
)

var (
	cfgFile = "../../../config-dev.yaml"
)

func TestVideoCoverSelection(t *testing.T) {
	var cfg struct {
		Vol conf.Vol
	}
	v := viper.New()
	v.SetConfigFile(cfgFile)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	// t.Log(v.AllSettings())
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}

	// 读入video_url.txt的内容给video_url
	r, err := os.OpenFile("video_url.txt", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	video_url, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	inst := CreateVisualService(cfg.Vol.Ak, cfg.Vol.Sk)
	form := url.Values{}

	form.Add("video_url", string(video_url))
	res, code, err := VideoCoverSelection(inst, form)
	if err != nil {
		t.Error(err)
	}

	t.Log(code, err)

	// res.Data.Results
	for idx, v := range res.Data.Results {
		img := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v.Data))
		buff, _ := io.ReadAll(img)
		fmt.Println(len(buff))
		// wtite file
		f, err := os.OpenFile(fmt.Sprintf("test_%d.jpg", idx), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			t.Error(err)
		}
		_, err = f.Write(buff)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("idx: %d, score: %f", idx, v.Score)
	}
}
