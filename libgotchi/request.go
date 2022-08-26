package libgotchi

import (
	"net/http"
	"strings"
	"time"
)

type ReqConf struct {
	Method  string
	Url     string
	Host    string
	Headers map[string]string
	Cookies map[string]string
	Data    []byte
}

func NewReqConf() ReqConf {
	var reqConf ReqConf
	reqConf.Method = "GET"
	reqConf.Url = ""
	reqConf.Host = ""
	reqConf.Headers = make(map[string]string)
	reqConf.Cookies = make(map[string]string)
	return reqConf
}

func SendRequest(reqConf *ReqConf, word string) Response {
	// *******************************************************************************
	// Replace EGG to word
	// *******************************************************************************
	reqConf.Method = strings.Replace(reqConf.Method, "EGG", word, -1)
	reqConf.Url = strings.Replace(reqConf.Url, "EGG", word, -1)
	reqConf.Host = strings.Replace(reqConf.Host, "EGG", word, -1)
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

	req, err := http.NewRequest(reqConf.Method, reqConf.Url, nil)
	if err != nil {
		panic(err)
	}

	// *******************************************************************************
	// Add custom headers
	// *******************************************************************************
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	for key, val := range reqConf.Headers {
		// Replace EGG to word
		key = strings.Replace(key, "EGG", word, -1)
		val = strings.Replace(val, "EGG", word, -1)
		req.Header.Add(key, val)
	}
	// Add custom cookies
	for key, val := range reqConf.Cookies {
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
