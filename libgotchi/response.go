package libgotchi

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Res struct {
	Body          string
	Config        Conf
	ContentLength int
	Duration      time.Duration
	Headers       map[string]string
	Status        string
	StatusCode    int
	Word          string
}

func NewResponse(resp *http.Response, req *Req, word string) Res {
	var r Res
	r.Config = req.Config
	r.ContentLength = int(resp.ContentLength)
	r.Duration = req.Duration
	r.Headers = make(map[string]string)
	r.Status = resp.Status
	r.StatusCode = resp.StatusCode
	r.Word = word

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		r.Body = string(body)
	} else {
		panic(err)
	}
	return r
}

func ErrorResponse(req *Req, word string) Res {
	var r Res
	r.Config = req.Config
	return r
}
