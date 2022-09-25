package libgotchi

import (
	"net/http"
	"strings"
	"time"
)

type Req struct {
	Config  Conf
	Cookies map[string]string
	Data    []byte
	Headers map[string]string
	Host    string
	Method  string
	Rate    time.Duration
	Url     string
}

func (r *Req) Send(word string) (Res, error) {
	r.Method = strings.Replace(r.Method, "EGG", word, -1)
	r.Url = strings.Replace(r.Url, "EGG", word, -1)
	r.Host = strings.Replace(r.Host, "EGG", word, -1)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       time.Duration(time.Duration(r.Config.Timeout) * time.Second),
		Transport: &http.Transport{
			// Proxy: nil,
			MaxConnsPerHost:     500,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     30 * time.Second,
			// DisableCompression:  true,
		},
	}

	req, err := http.NewRequest(r.Method, r.Url, nil)
	if err != nil {
		return ErrorResponse(r, word), err
	}

	// Fuzzing headers
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	for key, val := range r.Headers {
		// Replace EGG to word
		key = strings.Replace(key, "EGG", word, -1)
		val = strings.Replace(val, "EGG", word, -1)
		req.Header.Add(key, val)
	}
	// Fuzzing cookies
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

	resp, err := client.Do(req)
	if err != nil {
		return ErrorResponse(r, word), err
	}
	defer resp.Body.Close()

	response := NewResponse(resp, r, word)

	return response, nil
}

func NewReq(conf Conf) Req {
	var req Req
	req.Config = conf
	req.Cookies = make(map[string]string)
	req.Headers = make(map[string]string)
	req.Host = ""
	req.Method = conf.Method
	req.Rate = NewRate(conf.Rate)
	req.Url = conf.Url

	// Update headers
	if len(conf.Header) > 0 {
		headers := strings.Split(conf.Header, ";")
		for _, v := range headers {
			header := strings.Split(strings.TrimSpace(v), ":")
			key := header[0]
			val := header[1]
			req.Headers[key] = val
		}
	}
	// Update cookies
	if len(conf.Cookie) > 0 {
		cookies := strings.Split(conf.Cookie, ";")
		for _, v := range cookies {
			c := strings.Split(strings.TrimSpace(v), "=")
			key := c[0]
			val := c[1]
			req.Cookies[key] = val
		}
	}

	return req
}
