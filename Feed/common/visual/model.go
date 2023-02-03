package visual

import (
	volc "github.com/volcengine/volc-sdk-golang/base"
)

type VideoCoverSelectRequest struct {
	VideoUrl                    string   `json:"video_url"`
	VideoId                     string   `json:"video_id"`
	ImagesBase64                []string `json:"images_base64"`
	Enable                      bool     `json:"enable"`
	CutMethod                   string   `json:"cut_method"`
	Heights                     int      `json:"heights"`
	Widths                      int      `json:"widths"`
	UseRatio                    bool     `json:"use_ratio"`
	QualityMethod               string   `json:"quality_method"`
	PosterValidCheckerThreshold float64  `json:"poster_valid_checker_threshold"`
	ImageSelectorThreshold      float64  `json:"image_selector_threshold"`
}

type ImageInfoData struct {
	Score float64 `json:"score"`
	Data  string  `json:"data"`
}

type ImageInfoResult struct {
	Results []ImageInfoData `json:"results"`
}

type VideoCoverSelectResult struct {
	ResponseMetadata *volc.ResponseMetadata `json:",omitempty"`
	RequestId        string                 `json:"request_id"`
	TimeElapsed      string                 `json:"time_elapsed"`
	Code             int                    `json:"code"`
	Message          string                 `json:"message"`
	Data             *ImageInfoResult       `json:"data"`
}
