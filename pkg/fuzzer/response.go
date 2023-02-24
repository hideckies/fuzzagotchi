package fuzzer

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Body          string            `json:"body"`
	Config        Config            `json:"config"`
	ContentLength int               `json:"content_length"`
	Headers       map[string]string `json:"headers"`
	Delay         time.Duration     `json:"delay"`
	Status        string            `json:"status"`
	StatusCode    int               `json:"status_code"`
	Word          string            `json:"word"`
}

func NewResponse(resp *http.Response, req *Request, word string) Response {
	var newResp Response
	newResp.Config = req.Config
	newResp.ContentLength = int(resp.ContentLength)
	newResp.Headers = make(map[string]string)
	newResp.Delay = req.Delay
	newResp.Status = resp.Status
	newResp.StatusCode = resp.StatusCode
	newResp.Word = word

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		newResp.Body = string(body)
	} else {
		panic(err)
	}
	return newResp
}

// Error response
func errorResponse(req *Request, word string) Response {
	var resp Response
	resp.Config = req.Config
	return resp
}
