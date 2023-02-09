package video

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/volcengine/volc-sdk-golang/service/visual"
)

type Video interface {
	GetVideoCover(videoUrl string) ([]io.Reader, error)
}

type volcengine struct {
	accessKey string
	secretKey string
}

func NewVolcengineService(ak, sk string) Video {
	return &volcengine{
		accessKey: ak,
		secretKey: sk,
	}
}

func (s *volcengine) GetVideoCover(videoUrl string) ([]io.Reader, error) {
	form := url.Values{}
	form.Add("video_url", videoUrl)
	visual.DefaultInstance.Client.SetAccessKey(s.accessKey)
	visual.DefaultInstance.Client.SetSecretKey(s.secretKey)

	form.Add("video_url", videoUrl)

	resp, status, err := visual.DefaultInstance.VideoCoverSelection(form)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("volcengine video cover selection failed, status: %d", status)
	}
	covers := make([]io.Reader, len(resp.Data.Results))
	for _, v := range resp.Data.Results {
		img := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v.Data))
		covers = append(covers, img)
	}
	return covers, nil
}
