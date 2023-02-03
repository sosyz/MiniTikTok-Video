package visual

import (
	"encoding/json"
	"net/http"
	"net/url"

	volc "github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/visual"
)

type VisualInst struct {
	vis *visual.Visual
}

func CreateVisualService(accessKey, secretKey string) *VisualInst {
	instance := visual.NewInstance()
	instance.Client.SetAccessKey(accessKey)
	instance.Client.SetSecretKey(secretKey)

	instance.Client.ApiInfoList["VideoCoverSelection"] = &volc.ApiInfo{
		Method: http.MethodPost,
		Path:   "/",
		Query: url.Values{
			"Action":  []string{"VideoCoverSelection"},
			"Version": []string{"2020-08-26"},
		},
	}
	return &VisualInst{
		vis: instance,
	}
}

func VideoCoverSelection(instance *VisualInst, form url.Values) (*VideoCoverSelectResult, int, error) {
	resp := new(VideoCoverSelectResult)
	data, statusCode, err := instance.vis.Client.Post("VideoCoverSelection", nil, form)
	if err != nil {
		errMsg := err.Error()
		if errMsg[:3] != "api" {
			return nil, statusCode, err
		}
	}

	if err := json.Unmarshal(data, resp); err != nil {
		return nil, statusCode, err
	}
	return resp, statusCode, nil
}
