package fuzzer

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Client    *http.Client      `json:"client"`
	Config    Config            `json:"config"`
	Cookies   map[string]string `json:"cookies"`
	Data      []byte            `json:"data"`
	Delay     time.Duration     `json:"delay"`
	Headers   map[string]string `json:"headers"`
	Host      string            `json:"host"`
	Method    string            `json:"method"`
	PostData  io.Reader         `json:"post_data"`
	URL       string            `json:"url"`
	UserAgent string            `json:"user_agent"`
}

func NewRequest(conf Config) Request {
	var req Request
	req.Config = conf
	req.Cookies = make(map[string]string)
	req.Headers = make(map[string]string)
	req.Host = ""
	req.Method = conf.Method
	postdata := []byte(conf.PostData)
	req.PostData = bytes.NewReader(postdata)
	req.Delay = getDelay(conf.Delay)
	req.URL = conf.URL
	req.UserAgent = conf.UserAgent

	// // Update headers
	if len(conf.Header) > 0 {
		headers := strings.Split(conf.Header, ";")
		for _, v := range headers {
			header := strings.Split(strings.TrimSpace(v), ":")
			key := header[0]
			val := header[1]
			req.Headers[key] = val
		}
	}
	if _, ok := req.Headers["User-Agent"]; !ok {
		req.Headers["User-Agent"] = req.UserAgent
	}
	if _, ok := req.Headers["Accept-Language"]; !ok {
		req.Headers["Accept-Language"] = "en-US"
	}
	if _, ok := req.Headers["Connection"]; !ok {
		req.Headers["Connection"] = "Keep-Alive"
	}

	// // Update cookies
	if len(conf.Cookie) > 0 {
		cookies := strings.Split(conf.Cookie, ";")
		for _, v := range cookies {
			c := strings.Split(strings.TrimSpace(v), "=")
			key := c[0]
			val := c[1]
			req.Cookies[key] = val
		}
	}

	req.Client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       time.Duration(time.Duration(req.Config.Timeout) * time.Second),
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// Proxy: nil,
			MaxConnsPerHost:     500,
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 500,
			// IdleConnTimeout:     30 * time.Second,
			DialContext: (&net.Dialer{
				Timeout: time.Duration(time.Duration(req.Config.Timeout) * time.Second),
			}).DialContext,
			TLSHandshakeTimeout: time.Duration(time.Duration(req.Config.Timeout) * time.Second),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateOnceAsClient,
				ServerName:         "",
			},
		},
	}

	return req
}

func (req *Request) Send(word string) (Response, error) {
	// var resp Response

	req.Host = strings.ReplaceAll(req.Host, "EGG", word)
	req.Method = strings.ReplaceAll(req.Method, "EGG", word)
	postdata := []byte(strings.ReplaceAll(string(req.Config.PostData), "EGG", word))
	req.PostData = bytes.NewReader(postdata)
	if len(postdata) > 0 {
		req.Method = "POST"
	}
	req.URL = strings.ReplaceAll(req.URL, "EGG", word)

	newReq, err := http.NewRequest(req.Method, req.URL, req.PostData)
	if err != nil {
		return errorResponse(req, word), err
	}

	// Set Headers
	for key, val := range req.Headers {
		// Replace EGG to word
		key = strings.ReplaceAll(key, "EGG", word)
		val = strings.ReplaceAll(val, "EGG", word)
		newReq.Header.Set(key, val)
	}
	// Set Cookies
	for key, val := range req.Cookies {
		// Replace EGG to word
		key = strings.ReplaceAll(key, "EGG", word)
		val = strings.ReplaceAll(val, "EGG", word)
		cookie := &http.Cookie{
			Name:  key,
			Value: val,
		}
		newReq.AddCookie(cookie)
	}

	tmpResp, err := req.Client.Do(newReq)
	if err != nil {
		return errorResponse(req, word), err
	}
	defer tmpResp.Body.Close()

	resp := NewResponse(tmpResp, req, word)
	return resp, nil
}
