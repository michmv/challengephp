package request

import "net/http"

type Request struct {
	r *http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{r: r}
}

func (req *Request) Body() []byte {
	defer req.r.Body.Close()
	var body RequestBody
}
