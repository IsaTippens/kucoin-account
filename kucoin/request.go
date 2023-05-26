package kucoin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
)

type Request struct {
	fullUrl string
	fullPath string
	URL     string
	Path    string
	Method  string
	Headers http.Header
	Query   url.Values
	Body    []byte
}

type Response struct {
	Message []byte
	Status int
}

/*
GET/DELETE expects a map[string]string passed as params

GET https://example.com/?param1=2&param2=3
*/
func NewRequest(method string, path string, queries map[string]string, headers interface{}) Request {
	req := Request{
		Method:  method,
		Path:    path,
		Query:   make(url.Values),
		Headers: make(http.Header),
		Body:    make([]byte, 0),
	}

	req.buildParams(queries, headers)

	return req
}

func (req *Request) buildParams(queries map[string]string, headers interface{}) {
	if req.Method == http.MethodGet || req.Method == http.MethodDelete {
		if queries == nil {
			return
		}
		for k, v := range queries {
			req.Query.Set(k, v)
		}
	}

	if headers == nil {
		return
	}

	b, err := json.Marshal(headers)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Body = b
}

func (req *Request) GetFullUrl() string {
	if len(req.fullUrl) == 0 {
		req.generateFullUrl()
	}
	return req.fullUrl
}

func (req *Request) generateFullUrl() {
	req.fullUrl = fmt.Sprintf("%s%s", req.URL, req.Path)
	if len(req.Query) > 0 {
		if strings.Contains(req.fullUrl, "?") {
			req.fullUrl += "&" + req.Query.Encode()
		} else {
			req.fullUrl += "?" + req.Query.Encode()
		}
	}
}

func (req *Request) GetFullPath() string {
	if len(req.fullPath) == 0 {
		req.generateFullPath()
	}
	return req.fullPath
}


func (req *Request) generateFullPath() {
	req.fullPath = req.Path
	if len(req.Query) > 0 {
		if strings.Contains(req.fullPath, "?") {
			req.fullPath += "&" + req.Query.Encode()
		} else {
			req.fullPath += "?" + req.Query.Encode()
		}
	}
}


func (r *Request) fasthttpRequest() *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(r.Method)
	req.SetRequestURI(r.GetFullUrl())
	req.SetBody(r.Body)

	for key, values := range r.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return req
}
