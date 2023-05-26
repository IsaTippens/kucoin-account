package kucoin

import (
	"encoding/json"
)

const SuccessCode = "200000"

type KucoinApiConfig struct {
	ApiKey        string
	ApiSecret     string
	ApiPassPhrase string
}

type KucoinResponse struct {
	Message string	`json:"msg"`
	Code  string	`json:"code"`
	Data    json.RawMessage `json:"data"`
}

func ResponseFromJson(data []byte) (*KucoinResponse, error) {
	var kr KucoinResponse
	if err := json.Unmarshal(data, &kr); err != nil {
		return nil, err
	}
	return &kr, nil
}

func (kr KucoinResponse) Unmarshal(model interface{}) error {
	return json.Unmarshal(kr.Data, model)
}

