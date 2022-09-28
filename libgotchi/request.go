package libgotchi

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Req struct {
	Client   *http.Client
	Config   Conf
	Cookies  map[string]string
	Data     []byte
	Headers  map[string]string
	Host     string
	Method   string
	PostData io.Reader
	Rate     time.Duration
	Url      string
}

func (r *Req) Send(word string) (Res, error) {
	r.Host = strings.ReplaceAll(r.Host, "EGG", word)
	r.Method = strings.ReplaceAll(r.Method, "EGG", word)
	postdata := []byte(strings.ReplaceAll(string(r.Config.PostData), "EGG", word))
	r.PostData = bytes.NewReader(postdata)
	if len(postdata) > 0 {
		r.Method = "POST"
	}
	r.Url = strings.ReplaceAll(r.Url, "EGG", word)

	if _, ok := r.Headers["User-Agent"]; !ok {
		r.Headers["User-Agent"] = "Fuzzagotchi"
	}
	if _, ok := r.Headers["Connection"]; !ok {
		r.Headers["Connection"] = "Keep-Alive"
	}
	if _, ok := r.Headers["Accept-Language"]; !ok {
		r.Headers["Accept-Language"] = "en-US"
	}

	req, err := http.NewRequest(r.Method, r.Url, r.PostData)
	if err != nil {
		return ErrorResponse(r, word), err
	}

	// Headers
	for key, val := range r.Headers {
		// Replace EGG to word
		key = strings.ReplaceAll(key, "EGG", word)
		val = strings.ReplaceAll(val, "EGG", word)
		req.Header.Set(key, val)
	}
	// Cookies
	for key, val := range r.Cookies {
		// Replace EGG to word
		key = strings.ReplaceAll(key, "EGG", word)
		val = strings.ReplaceAll(val, "EGG", word)
		cookie := &http.Cookie{
			Name:  key,
			Value: val,
		}
		req.AddCookie(cookie)
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return ErrorResponse(r, word), err
	}
	defer resp.Body.Close()

	response := NewResponse(resp, r, word)

	return response, nil
}

func NewReq(conf Conf) Req {
	var r Req
	r.Config = conf
	r.Cookies = make(map[string]string)
	r.Headers = make(map[string]string)
	r.Host = ""
	r.Method = conf.Method
	postdata := []byte(conf.PostData)
	r.PostData = bytes.NewReader(postdata)
	r.Rate = NewRate(conf.Rate)
	r.Url = conf.Url

	// Update headers
	if len(conf.Header) > 0 {
		headers := strings.Split(conf.Header, ";")
		for _, v := range headers {
			header := strings.Split(strings.TrimSpace(v), ":")
			key := header[0]
			val := header[1]
			r.Headers[key] = val
		}
	}
	// Update cookies
	if len(conf.Cookie) > 0 {
		cookies := strings.Split(conf.Cookie, ";")
		for _, v := range cookies {
			c := strings.Split(strings.TrimSpace(v), "=")
			key := c[0]
			val := c[1]
			r.Cookies[key] = val
		}
	}

	r.Client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		Timeout:       time.Duration(time.Duration(r.Config.Timeout) * time.Second),
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			// Proxy: nil,
			MaxConnsPerHost:     500,
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 500,
			// IdleConnTimeout:     30 * time.Second,
			DialContext: (&net.Dialer{
				Timeout: time.Duration(time.Duration(r.Config.Timeout) * time.Second),
			}).DialContext,
			TLSHandshakeTimeout: time.Duration(time.Duration(r.Config.Timeout) * time.Second),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateOnceAsClient,
				ServerName:         "",
			},
		},
	}

	return r
}
