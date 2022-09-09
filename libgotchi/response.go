package libgotchi

import (
	"io/ioutil"
	"net/http"
)

type Res struct {
	Body          string
	ContentLength int
	Headers       map[string]string
	Status        string
	StatusCode    int
}

func NewResponse(resp *http.Response) Res {
	var r Res
	r.ContentLength = int(resp.ContentLength)
	r.Headers = make(map[string]string)
	r.Status = resp.Status
	r.StatusCode = resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		r.Body = string(body)
	} else {
		panic(err)
	}
	return r
}
