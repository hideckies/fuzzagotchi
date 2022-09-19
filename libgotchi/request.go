package libgotchi

import (
	"net/http"
	"strings"
	"time"
)

type Req struct {
	Config   Conf
	Cookies  map[string]string
	Data     []byte
	Duration time.Duration
	Headers  map[string]string
	Host     string
	Method   string
	Url      string
}

func (r *Req) Send(word string) (Res, error) {
	r.Method = strings.Replace(r.Method, "EGG", word, -1)
	r.Url = strings.Replace(r.Url, "EGG", word, -1)
	r.Host = strings.Replace(r.Host, "EGG", word, -1)

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
		return ErrorResponse(r, word), err
	}

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
	req.Duration = NewDuration(conf.TimeDelay)
	req.Headers = make(map[string]string)
	req.Host = ""
	req.Method = conf.Method
	req.Url = conf.Url
	return req
}
