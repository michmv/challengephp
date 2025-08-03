package stats_handler

import (
	"challengephp/lib"
	"challengephp/src/response"
	"challengephp/src/services/repository"
	"challengephp/src/services/response_service"
	"fmt"
	"net/http"
	"time"
)

type Services interface {
	ResponseFactory() response_service.ResponseFactory
	Repository() repository.Repository
	Request() *http.Request
}

func Handler(serv Services) (response.Response, lib.Error) {
	typeName, from, to, err2 := getQueryParam(serv.Request())
	if err2 != nil {
		return serv.ResponseFactory().Error(404, err2.Error()), nil
	}

	eventType, ok, err := serv.Repository().FindEventType(typeName)
	if err != nil {
		return nil, err.Tap()
	}
	if !ok {
		return serv.ResponseFactory().Error(404, "не найден тип"), nil
	}

	pages, err := serv.Repository().GetStatsPages(eventType.Id, from, to)
	if err != nil {
		return nil, err.Tap()
	}
	usersUnique, err := serv.Repository().GetStatUniqueUsers(eventType.Id, from, to)
	if err != nil {
		return nil, err.Tap()
	}

	total, err := serv.Repository().GetStatTotal(eventType.Id, from, to)
	if err != nil {
		return nil, err.Tap()
	}

	res := Response{
		TotalEvents: total,
		UniqueUsers: usersUnique,
		TopPages:    pages,
	}
	return serv.ResponseFactory().Json(res), nil
}

type Response struct {
	TotalEvents int64            `json:"total_events"`
	UniqueUsers int64            `json:"unique_users"`
	TopPages    map[string]int64 `json:"top_pages"`
}

func getQueryParam(r *http.Request) (string, string, string, error) {
	fromStr := r.URL.Query().Get("from")
	if fromStr != "" {
		_, err := time.Parse("2006-01-02 15:04:05", fromStr)
		if err != nil {
			return "", "", "", fmt.Errorf("Неудалось распарсить from")
		}
	}
	toStr := r.URL.Query().Get("to")
	if toStr != "" {
		_, err := time.Parse("2006-01-02 15:04:05", toStr)
		if err != nil {
			return "", "", "", fmt.Errorf("Неудалось распарсить to")
		}
	}
	typeStr := r.URL.Query().Get("type")
	return typeStr, fromStr, toStr, nil
}
