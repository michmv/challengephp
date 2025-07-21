package ping_handler

import (
	"challengephp/lib"
	"challengephp/src/response"
	"challengephp/src/services/response_service"
)

type Services interface {
	ResponseFactory() response_service.ResponseFactory
}

func Handler(serv Services) (response.Response, lib.Error) {
	return serv.ResponseFactory().String("ok"), nil
}
