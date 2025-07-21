package response_service

import (
	"challengephp/lib"
	"encoding/json"
	"net/http"
)

type ResponseFactory struct {
}

func NewResponseFactory() ResponseFactory {
	return ResponseFactory{}
}

func (it ResponseFactory) String(str string) StringResponse {
	return NewStringResponse(str)
}

func (it ResponseFactory) Json(data any) JsonResponse {
	return NewJsonResponse(data)
}

func (it ResponseFactory) Error(code int, str string) ErrorResponse {
	return NewErrorResponse(code, str)
}

type StringResponse struct {
	str string
}

func NewStringResponse(str string) StringResponse {
	return StringResponse{str}
}

func (it StringResponse) Render(w http.ResponseWriter) lib.Error {
	_, err := w.Write([]byte(it.str))
	if err != nil {
		return lib.Err(err)
	}
	return nil
}

type JsonResponse struct {
	data interface{}
}

func NewJsonResponse(data interface{}) JsonResponse {
	return JsonResponse{data}
}

func (it JsonResponse) Render(w http.ResponseWriter) lib.Error {
	jsonData, err := json.Marshal(it.data)
	if err != nil {
		return lib.Err(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		return lib.Err(err)
	}
	return nil
}

type ErrorResponse struct {
	code int
	str  string
}

func NewErrorResponse(code int, str string) ErrorResponse {
	return ErrorResponse{code, str}
}

func (it ErrorResponse) Render(w http.ResponseWriter) lib.Error {
	e := Error{Code: it.code, Message: it.str, Error: http.StatusText(it.code)}
	jsonData, err := json.Marshal(e)
	if err != nil {
		return lib.Err(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(it.code)
	_, err = w.Write(jsonData)
	if err != nil {
		return lib.Err(err)
	}
	return nil
}

type Error struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
