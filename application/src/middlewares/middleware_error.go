package middlewares

import (
	"challengephp/lib"
	"challengephp/src/response"
	"challengephp/src/services/response_service"
	"fmt"
	"net/http"
	"runtime/debug"
)

type ErrorServices interface {
	Debug() bool
	Logger() lib.Logger
	ResponseFactory() response_service.ResponseFactory
	ResponseWriter() http.ResponseWriter
}

func ErrorMiddleware(serv ErrorServices, next lib.ActionFunc[response.Response]) (response response.Response, err lib.Error) {
	defer func() {
		if err_ := recover(); err_ != nil {
			err = lib.ErrPanic(fmt.Errorf("panic: %v", err_), debug.Stack())
			err = errorsHandler(serv.ResponseFactory(), serv.ResponseWriter(), err, serv.Debug(), serv.Logger())
			response = nil
		}
	}()

	response, err = next()
	if err != nil {
		err = errorsHandler(serv.ResponseFactory(), serv.ResponseWriter(), err, serv.Debug(), serv.Logger())
		return nil, err // ошибку отобразили
	}

	return nil, nil
}

func errorsHandler(responseFactory response_service.ResponseFactory, w http.ResponseWriter, errIn lib.Error, debugOn bool, logger lib.Logger) lib.Error {
	if logger != nil {
		logger.Error(errIn.Error())
	}

	textError := "***"
	if debugOn {
		textError = errIn.Error()
	}

	res := responseFactory.Error(http.StatusInternalServerError, textError)
	err := res.Render(w)
	if err != nil {
		return err.Add("ошибка при создании страницы отображения")
	}
	return nil
}
