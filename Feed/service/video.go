package service

import (
	"encoding/base64"
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/sosyz/mini_tiktok_feed/Feed/common/visual"
)

type VideoService interface {
	GetVideoCover(videoUrl string) (io.Reader, error)
}

type volcengine struct {
	instance *visual.VisualInst
}

func NewVolcengineService(ak, sk string) VideoService {
	return &volcengine{
		instance: visual.CreateVisualService(ak, sk),
	}
}

func (s *volcengine) GetVideoCover(videoUrl string) (io.Reader, error) {
	form := url.Values{}

	form.Add("video_url", videoUrl)
	res, _, err := visual.VideoCoverSelection(s.instance, form)
	if err != nil {
		return nil, err
	}

	for _, v := range res.Data.Results {
		img := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v.Data))
		return img, nil
	}
	return nil, errors.New("no cover found")
}
