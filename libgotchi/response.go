package libgotchi

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	Body       string
	Headers    map[string]string
	Status     string
	StatusCode int
}

func NewResponse(resp *http.Response) Response {
	var response Response
	response.Headers = make(map[string]string)
	response.Status = resp.Status
	response.StatusCode = resp.StatusCode

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		response.Body = string(body)
	} else {
		panic(err)
	}
	return response
}
