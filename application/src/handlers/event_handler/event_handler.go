package event_handler

import (
	"challengephp/lib"
	"challengephp/src/response"
	"challengephp/src/services/repository"
	"challengephp/src/services/response_service"
	"challengephp/src/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Services interface {
	ResponseFactory() response_service.ResponseFactory
	Repository() repository.Repository
	Request() *http.Request
}

func HandlerGet(serv Services) (response.Response, lib.Error) {
	page, limit, err2 := getRequestQuery(serv.Request())
	if err2 != nil {
		return serv.ResponseFactory().Error(400, err2.Error()), nil
	}
	events, err := serv.Repository().List(page, limit)
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
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
	return serv.ResponseFactory().Json(res), nil
}

func HandlerPost(serv Services) (response.Response, lib.Error) {
	data, err2 := parseJson(serv.Request())
	if err2 != nil {
		return serv.ResponseFactory().Error(400, "Не получилось распарсить запрос"), nil
	}

	user, err := serv.Repository().FindOrCreateUser(data.UserId)
	if err != nil {
		return nil, err.Tap()
	}

	eventType, err := serv.Repository().FindOrCreateEvent(data.EventType)
	if err != nil {
		return nil, err.Tap()
	}

	event := types.Event{
		Id:        0,
		Timestamp: data.Timestamp,
		Metadata:  data.Metadata,
		UserId:    user.Id,
		TypeId:    eventType.Id,
	}

	id, err := serv.Repository().CreateEvent(event)
	if err != nil {
		return nil, err.Tap()
	}
	event.Id = id

	return serv.ResponseFactory().Json(event), nil
}

func getRequestQuery(r *http.Request) (page, limit int64, err error) {
	page_ := r.URL.Query().Get("page")
	if page_ == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(page_, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("page должно быть целым числом")
		}
	}
	limit_ := r.URL.Query().Get("limit")
	if limit_ == "" {
		limit = 100
	} else {
		limit, err = strconv.ParseInt(limit_, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("limit должно быть целым числом")
		}
	}
	return page, limit, nil
}

func parseJson(r *http.Request) (RequestBody, error) {
	defer r.Body.Close()
	var body RequestBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}

type RequestBody struct {
	UserId    int64                  `json:"user_id"`
	EventType string                 `json:"event_type"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type Response struct {
	Data  []types.Event2 `json:"data"`
	Query QueryResponse  `json:"query"`
}

type QueryResponse struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
	Total int64 `json:"total"`
}
