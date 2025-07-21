package user_handler

import (
	"challengephp/lib"
	"challengephp/src/response"
	"challengephp/src/services/repository"
	"challengephp/src/services/response_service"
	"challengephp/src/types"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Services interface {
	ResponseFactory() response_service.ResponseFactory
	Repository() repository.Repository
	Request() *http.Request
}

func Handler(serv Services) (response.Response, lib.Error) {
	userId, err2 := getRequestQuery(serv.Request())
	if err2 != nil {
		return serv.ResponseFactory().Error(404, "Страница не найдена"), nil
	}

	user, ok, err := serv.Repository().FindUser(userId)
	if err != nil {
		return nil, err.Tap()
	}
	if !ok {
		return serv.ResponseFactory().Error(404, "Пользователь не найден"), nil
	}
	events, err := serv.Repository().ListForUser(user.Id)
	if err != nil {
		return nil, err.Tap()
	}
	total, err := serv.Repository().EventsTotal()
	if err != nil {
		return nil, err.Tap()
	}
	res := Response{
		Data: events,
		Query: QueryResponse{
			Total: total,
			Limit: 1000,
		},
	}
	return serv.ResponseFactory().Json(res), nil
}

func getRequestQuery(r *http.Request) (int64, error) {
	userId_ := chi.URLParam(r, "userId")
	userId, err := strconv.ParseInt(userId_, 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

type Response struct {
	Data  []types.Event2 `json:"data"`
	Query QueryResponse  `json:"query"`
}

type QueryResponse struct {
	Limit int64 `json:"limit"`
	Total int64 `json:"total"`
}
