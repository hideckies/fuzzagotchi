package fuzzer

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	Client         *http.Client      `json:"client"`
	Config         Config            `json:"config"`
	Cookies        map[string]string `json:"cookies"`
	Data           []byte            `json:"data"`
	Delay          time.Duration     `json:"delay"`
	FollowRedirect bool              `json:"follow_redirect"`
	Headers        map[string]string `json:"headers"`
	Host           string            `json:"host"`
	Method         string            `json:"method"`
	PostData       io.Reader         `json:"post_data"`
	URL            string            `json:"url"`
	UserAgent      string            `json:"user_agent"`
}

// Initialize Request
func NewRequest(conf Config) Request {
	proxyURL := http.ProxyFromEnvironment
	customProxy := ""
	if conf.Proxy != "" {
		customProxy = conf.Proxy
		p, err := url.Parse(customProxy)
		if err == nil {
			proxyURL = http.ProxyURL(p)
		}
	}

	var req Request
	req.Config = conf
	req.Cookies = make(map[string]string)
	req.Delay = getDelay(conf.Delay)
	req.FollowRedirect = conf.FollowRedirect
	req.Headers = make(map[string]string)
	req.Host = conf.Host
	req.Method = conf.Method
	postdata := []byte(conf.PostData)
	req.PostData = bytes.NewReader(postdata)
	req.URL = conf.URL
	req.UserAgent = conf.UserAgent

	// Set headers
	if len(conf.Header) > 0 {
		// If Cookie
		if strings.Contains(conf.Header, "Cookie") {
			cookie := strings.Replace(strings.TrimSpace(conf.Header), "Cookie:", "", -1)
			cookieVals := strings.Split(cookie, ";")
			for _, cookieVal := range cookieVals {
				c := strings.Split(strings.TrimSpace(cookieVal), "=")
				key := c[0]
				val := c[1]
				req.Cookies[key] = val
			}
		} else {
			// If common headers
			headers := strings.Split(conf.Header, ";")
			for _, v := range headers {
				header := strings.Split(strings.TrimSpace(v), ":")
				key := header[0]
				val := header[1]
				req.Headers[key] = val
			}
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

	// Set cookies
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
			ForceAttemptHTTP2:   true,
			Proxy:               proxyURL,
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

	// if req.Config.FollowRedirect {
	// 	req.Client.CheckRedirect = nil
	// }

	return req
}

// Sent request
func (req *Request) Send(word string) (Response, error) {
	newReqHeaders := req.Headers
	newReqMethod := strings.ReplaceAll(req.Method, "EGG", word)
	newReqCookies := req.Cookies
	tmpPostData := []byte(strings.ReplaceAll(string(req.Config.PostData), "EGG", word))
	newReqPostData := bytes.NewReader(tmpPostData)
	if len(tmpPostData) > 0 {
		newReqMethod = "POST"
	}
	newReqURL := strings.ReplaceAll(req.URL, "EGG", word)

	// Initialize a new http request
	newReq, err := http.NewRequest(newReqMethod, newReqURL, newReqPostData)
	if err != nil {
		return errorResponse(req, word), err
	}

	// Update Headers (also Host)
	for key, val := range newReqHeaders {
		// Replace EGG to word
		key = strings.TrimSpace(strings.ReplaceAll(key, "EGG", word))
		val = strings.TrimSpace(strings.ReplaceAll(val, "EGG", word))
		newReq.Header.Set(key, val)

		// If "Host" header exists, update newReq.Host
		if key == "Host" {
			newReq.Host = val
		}
	}
	if newReq.Header.Get("Content-Type") == "" {
		newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	// Update Cookies
	for key, val := range newReqCookies {
		// Replace EGG to word
		key = strings.TrimSpace(strings.ReplaceAll(key, "EGG", word))
		val = strings.TrimSpace(strings.ReplaceAll(val, "EGG", word))
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

	// Process redirects
	if tmpResp.StatusCode == http.StatusMovedPermanently || tmpResp.StatusCode == http.StatusFound || tmpResp.StatusCode == http.StatusSeeOther || tmpResp.StatusCode == http.StatusTemporaryRedirect {
		redirectUrl, err := tmpResp.Location()
		if err != nil {
			return errorResponse(req, word), err
		}
		newReq, err = http.NewRequest("GET", redirectUrl.String(), nil)
		if err != nil {
			return errorResponse(req, word), err
		}
		var redirectResp *http.Response
		redirectResp, err = req.Client.Do(newReq)
		if err != nil {
			return errorResponse(req, word), err
		}
		defer redirectResp.Body.Close()

		resp := NewResponse(tmpResp, req, word, getPath(newReqURL), redirectResp)
		return resp, nil
	}

	resp := NewResponse(tmpResp, req, word, getPath(newReqURL), nil)
	return resp, nil
}

// Get path from URL
func getPath(url string) string {
	urlSplit := strings.Split(url, "/")
	return "/" + strings.Join(urlSplit[3:], "/")
}
