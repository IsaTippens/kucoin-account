package kucoin

import (
	"bytes"
	"errors"
	"fmt"
)

type Client struct {
	ApiKey        string
	ApiSecret     string
	ApiPassphrase string
	signer        *KucoinSigner
	httpclient   *HttpClient
}

func NewClient(apiKey string, apiSecret string, apiPassphrase string) *Client {
	return &Client{
		ApiKey:        apiKey,
		ApiSecret:     apiSecret,
		ApiPassphrase: apiPassphrase,
		signer:        NewKcSigner(apiKey, apiSecret, apiPassphrase),
		httpclient:    NewHttpClient(),
	}
}

func (c Client) Do(req *Request) (*KucoinResponse, error) {
	req.Headers.Set("Content-Type", "application/json")
	var b bytes.Buffer
	b.WriteString(req.Method)
	b.WriteString(req.GetFullPath())
	b.Write(req.Body)

	sign_headers := c.signer.Headers(b.String())
	for k,v := range sign_headers {
		req.Headers.Set(k, v)
	}
	
	resp, err := c.httpclient.CallRequest(req)
	
	if err != nil {
		return nil, err
	}
	
	kr, err := ResponseFromJson(resp.Message)
	if err != nil {
		return nil, err
	}

	if kr.Code != SuccessCode {
		m := fmt.Sprintf("Error: Code %s Message %s", kr.Code, kr.Message)
		return kr, errors.New(m)
	}

	return kr, nil
}
