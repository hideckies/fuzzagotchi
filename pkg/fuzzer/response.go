package fuzzer

import (
	"io"
	"net/http"
	"time"
)

type Response struct {
	Body          []byte            `json:"body"`
	Config        Config            `json:"config"`
	ContentLength int               `json:"content_length"`
	Delay         time.Duration     `json:"delay"`
	Headers       map[string]string `json:"headers"`
	Path          string            `json:"path"`
	RedirectPath  string            `json:"redirect_path"`
	Status        string            `json:"status"`
	StatusCode    int               `json:"status_code"`
	Word          string            `json:"word"`
}

func NewResponse(resp *http.Response, req *Request, word string, redirectUrl string) Response {
	var newResp Response
	newResp.Body = make([]byte, 0)
	newResp.Config = req.Config
	newResp.ContentLength = int(resp.ContentLength)
	newResp.Delay = req.Delay
	newResp.Headers = make(map[string]string)
	newResp.Path = resp.Request.URL.Path
	newResp.RedirectPath = redirectUrl
	newResp.Status = resp.Status
	newResp.StatusCode = resp.StatusCode
	newResp.Word = word

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		newResp.Body = make([]byte, 0)
	} else {
		newResp.Body = body
	}

	if newResp.ContentLength < 0 {
		newResp.ContentLength = len(body)
	}

	return newResp
}

// Error response
func errorResponse(req *Request, word string) Response {
	var resp Response
	resp.Config = req.Config
	return resp
}
