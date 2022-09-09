package libgotchi

import (
	"net/http"
	"strings"
	"time"
)

type Req struct {
	Method  string
	Url     string
	Host    string
	Headers map[string]string
	Cookies map[string]string
	Data    []byte
}

func (r *Req) Send(word string) Res {
	// *******************************************************************************
	// Replace EGG to word
	// *******************************************************************************
	r.Method = strings.Replace(r.Method, "EGG", word, -1)
	r.Url = strings.Replace(r.Url, "EGG", word, -1)
	r.Host = strings.Replace(r.Host, "EGG", word, -1)
	// *******************************************************************************

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
	}

	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		panic(err)
	}

	// *******************************************************************************
	// Add custom headers
	// *******************************************************************************
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	for key, val := range r.Headers {
		// Replace EGG to word
		key = strings.Replace(key, "EGG", word, -1)
		val = strings.Replace(val, "EGG", word, -1)
		req.Header.Add(key, val)
	}
	// Add custom cookies
	for key, val := range r.Cookies {
		// Replace EGG to word
		key = strings.Replace(key, "EGG", word, -1)
		val = strings.Replace(val, "EGG", word, -1)
		cookie := &http.Cookie{
			Name:  key,
			Value: val,
		}
		req.AddCookie(cookie)
	}
	// *******************************************************************************

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response := NewResponse(resp)

	return response
}

func NewReq() Req {
	var req Req
	req.Method = "GET"
	req.Url = ""
	req.Host = ""
	req.Headers = make(map[string]string)
	req.Cookies = make(map[string]string)
	return req
}
