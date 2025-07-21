package routes

import (
	"challengephp/lib"
	"challengephp/src/di"
	"challengephp/src/handlers/event_handler"
	ping2 "challengephp/src/handlers/ping_handler"
	"challengephp/src/handlers/stats_handler"
	"challengephp/src/handlers/user_handler"
	"challengephp/src/middlewares"
	"challengephp/src/response"
	"challengephp/src/storage"
)

func RegistrationRoutes(rm *lib.RouterManager[*storage.Storage, *di.Container, response.Response]) {
	rm.Use(errorMiddleware, responseMiddleware)

	rm.Get("/ping", ping)
	rm.Post("/events", eventPost)
	rm.Get("/events", eventGet)
	rm.Get("/users/{userId}/events", user)
	rm.Get("/stats", stats)
}

func stats(di *di.Container) (response.Response, lib.Error) {
	return stats_handler.Handler(di)
}

func user(di *di.Container) (response.Response, lib.Error) {
	return user_handler.Handler(di)
}

func eventPost(di *di.Container) (response.Response, lib.Error) {
	return event_handler.HandlerPost(di)
}

func eventGet(di *di.Container) (response.Response, lib.Error) {
	return event_handler.HandlerGet(di)
}

func ping(di *di.Container) (response.Response, lib.Error) {
	return ping2.Handler(di)
}

func errorMiddleware(di *di.Container, next lib.ActionFunc[response.Response]) (response.Response, lib.Error) {
	return middlewares.ErrorMiddleware(di, next)
}

func responseMiddleware(di *di.Container, next lib.ActionFunc[response.Response]) (response.Response, lib.Error) {
	return middlewares.ResponseMiddleware(di, next)
}
