package libgotchi

import (
	"net/http"
	"time"
)

type ReqConf struct {
	Method  string
	Host    string
	Url     string
	Headers map[string]string
	Data    []byte
}

func NewReqConf() ReqConf {
	var reqConf ReqConf
	reqConf.Headers = make(map[string]string)
	reqConf.Method = "GET"
	reqConf.Url = ""
	return reqConf
}

func SendRequest(reqConf *ReqConf) Response {
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
	req.Header.Add("If-None-Match", `W/"wyzzy"`)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	response := NewResponse(resp)

	return response
}
