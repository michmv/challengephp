package middlewares

import (
	"challengephp/lib"
	"challengephp/src/response"
	"net/http"
)

type ResponseService interface {
	ResponseWriter() http.ResponseWriter
}

func ResponseMiddleware(serv ResponseService, next lib.ActionFunc[response.Response]) (response.Response, lib.Error) {
	res, err := next()
	if err != nil {
		return nil, err
	}
	err = res.Render(serv.ResponseWriter())
	if err != nil {
		return nil, err.Tap()
	}
	return nil, nil
}
