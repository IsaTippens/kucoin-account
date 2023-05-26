package kucoin

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

const KucoinURL = "https://api.kucoin.com"

type HttpClient struct {
	client fasthttp.Client
	BaseURL string
}

func NewHttpClient() *HttpClient {
	c := HttpClient{
		client: fasthttp.Client{},
	}
	c.BaseURL = KucoinURL

	return &c
}



func (h *HttpClient) CallRequest(r *Request) (*Response, error) {
	r.URL = h.BaseURL

	req := r.fasthttpRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := h.client.Do(req, resp); err != nil {
		fmt.Println(err)
		return nil, err
	}

	var hresp Response
	hresp.Message = resp.Body()
	hresp.Status = resp.StatusCode()
	return &hresp, nil
}
