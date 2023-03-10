package fuzzer

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Body          []byte            `json:"body"`
	Config        Config            `json:"config"`
	Content       string            `jsong:"content"`
	ContentLength int               `json:"content_length"`
	ContentWords  int               `jsong:"content_words"`
	Delay         time.Duration     `json:"delay"`
	Headers       map[string]string `json:"headers"`
	Path          string            `json:"path"`
	RedirectPath  string            `json:"redirect_path"`
	Status        string            `json:"status"`
	StatusCode    int               `json:"status_code"`
	Word          string            `json:"word"`
}

func NewResponse(resp *http.Response, req *Request, word string, reqPath string, redirectResp *http.Response) Response {
	var newResp Response
	newResp.Body = make([]byte, 0)
	newResp.Config = req.Config
	newResp.Content = ""
	newResp.ContentLength = int(resp.ContentLength)
	newResp.ContentWords = 0
	newResp.Delay = req.Delay
	newResp.Headers = make(map[string]string)
	newResp.Path = reqPath
	newResp.RedirectPath = ""
	newResp.Status = resp.Status
	newResp.StatusCode = resp.StatusCode
	newResp.Word = word

	reader := resp.Body

	if redirectResp != nil {
		newResp.RedirectPath = redirectResp.Request.URL.Path
		if req.Config.FollowRedirect {
			newResp.Status = redirectResp.Status
			newResp.StatusCode = redirectResp.StatusCode
			reader = redirectResp.Body
		}
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		newResp.Body = make([]byte, 0)
	} else {
		newResp.Body = body
		newResp.Content = string(body[:])
		words := strings.Split(newResp.Content, " ")
		newResp.ContentWords = len(words)
	}

	// Update content length
	newResp.ContentLength = len(body)

	return newResp
}

// Error response
func errorResponse(req *Request, word string) Response {
	var resp Response
	resp.Config = req.Config
	return resp
}
