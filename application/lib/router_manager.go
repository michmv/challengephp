package lib

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"runtime/debug"
)

type MiddlewareFunc[C, R any] func(C, ActionFunc[R]) (R, Error)

type HandlerFunc[C, R any] func(C) (R, Error)

type ActionFunc[R any] func() (R, Error)

type RouterManager[S, C, R any] struct {
	mux              *chi.Mux
	storage          S
	fnMakeContainer  func(S, http.ResponseWriter, *http.Request) C
	fnCloseContainer func(C)
	debug            bool
	logger           Logger
	middlewares      []MiddlewareFunc[C, R]
}

func NewRouterManager[S, C, R any](
	mux *chi.Mux,
	storage S,
	fnMakeContainer func(S, http.ResponseWriter, *http.Request) C,
	fnCloseContainer func(C),
	debug bool,
	logger Logger,
) *RouterManager[S, C, R] {
	return &RouterManager[S, C, R]{
		mux:              mux,
		storage:          storage,
		fnMakeContainer:  fnMakeContainer,
		fnCloseContainer: fnCloseContainer,
		debug:            debug,
		logger:           logger,
		middlewares:      make([]MiddlewareFunc[C, R], 0),
	}
}

func (it *RouterManager[S, C, R]) Use(middlewares ...MiddlewareFunc[C, R]) {
	it.middlewares = append(it.middlewares, middlewares...)
}

func (it *RouterManager[S, C, R]) Group(relativePath string) *Group[S, C, R] {
	return NewGroup(it, relativePath, it.middlewares...)
}

func (it *RouterManager[S, C, R]) NotFound(handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.NotFound(it.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *RouterManager[S, C, R]) NoMethod(handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.MethodNotAllowed(it.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *RouterManager[S, C, R]) Get(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.Get(relativePath, it.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *RouterManager[S, C, R]) Post(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.Post(relativePath, it.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *RouterManager[S, C, R]) List(methods []string, relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	for _, m := range methods {
		it.mux.Method(m, relativePath, it.makeHandler(handler, append(it.middlewares, middlewares...)))
	}
}

func (it *RouterManager[S, C, R]) Common(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.List([]string{http.MethodGet, http.MethodPost, http.MethodOptions}, relativePath, handler, middlewares...)
}

func (it *RouterManager[S, C, R]) All(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.List([]string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	},
		relativePath,
		handler,
		middlewares...,
	)
}

type Group[S, C, R any] struct {
	routerManager *RouterManager[S, C, R]
	mux           *chi.Mux
	middlewares   []MiddlewareFunc[C, R]
}

func NewGroup[S, C, R any](routerManager *RouterManager[S, C, R], relativePath string, middlewares ...MiddlewareFunc[C, R]) *Group[S, C, R] {
	mux := chi.NewRouter()
	routerManager.mux.Mount(relativePath, mux)
	return &Group[S, C, R]{
		routerManager: routerManager,
		mux:           mux,
		middlewares:   middlewares,
	}
}

func (it *Group[S, C, R]) Use(middlewares ...MiddlewareFunc[C, R]) {
	it.middlewares = append(it.middlewares, middlewares...)
}

func (it *Group[S, C, R]) Group(relativePath string, middlewares ...MiddlewareFunc[C, R]) *Group[S, C, R] {
	mux := chi.NewRouter()
	it.mux.Mount(relativePath, mux)
	return &Group[S, C, R]{
		routerManager: it.routerManager,
		mux:           mux,
		middlewares:   append(it.middlewares, middlewares...),
	}
}

func (it *Group[S, C, R]) NotFound(handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.NotFound(it.routerManager.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *Group[S, C, R]) NoMethod(handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.MethodNotAllowed(it.routerManager.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *Group[S, C, R]) Get(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.Get(relativePath, it.routerManager.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *Group[S, C, R]) Post(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.mux.Post(relativePath, it.routerManager.makeHandler(handler, append(it.middlewares, middlewares...)))
}

func (it *Group[S, C, R]) List(methods []string, relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	for _, m := range methods {
		it.mux.Method(m, relativePath, it.routerManager.makeHandler(handler, append(it.middlewares, middlewares...)))
	}
}

func (it *Group[S, C, R]) Common(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.List([]string{http.MethodGet, http.MethodPost, http.MethodOptions}, relativePath, handler, middlewares...)
}

func (it *Group[S, C, R]) All(relativePath string, handler HandlerFunc[C, R], middlewares ...MiddlewareFunc[C, R]) {
	it.List([]string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	},
		relativePath,
		handler,
		middlewares...,
	)
}

func (it *RouterManager[S, C, R]) makeHandler(handler HandlerFunc[C, R], middlewares []MiddlewareFunc[C, R]) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err Error
		defer func() {
			if err_ := recover(); err_ != nil {
				err = ErrPanic(fmt.Errorf("panic: %v", err_), debug.Stack())
				it.errorsHandler(writer, err, it.debug, it.logger)
			}
		}()

		container := it.fnMakeContainer(it.storage, writer, request)
		defer it.fnCloseContainer(container)

		err = it.pipeline(handler, middlewares, container)
		if err != nil {
			it.errorsHandler(writer, err, it.debug, it.logger)
		}
	}
}

func (it *RouterManager[S, C, R]) pipeline(handler HandlerFunc[C, R], middlewares []MiddlewareFunc[C, R], context C) (err Error) {
	fnc := func() (R, Error) {
		return handler(context)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		fnc = it.createMiddleware(fnc, middlewares[i], context)
	}
	_, err = fnc()
	if err != nil {
		_ = err.Tap()
	}
	return
}

func (it *RouterManager[S, C, R]) createMiddleware(next ActionFunc[R], middleware MiddlewareFunc[C, R], container C) ActionFunc[R] {
	return func() (R, Error) {
		return middleware(container, next)
	}
}

func (it *RouterManager[S, C, R]) errorsHandler(writer http.ResponseWriter, err Error, debugOn bool, logger Logger) {
	if logger != nil {
		logger.Error(err.Error())
	}

	message := ""
	if debugOn {
		message = fmt.Sprintf("<pre>%s</pre>", err.Error())
	}
	message = fmt.Sprintf("<h1>%d %s</h1>\n%s\n", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), message)
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusInternalServerError)
	_, err_ := writer.Write([]byte(message))
	if err_ != nil {
		panic(ErrPanic(err_, debug.Stack()))
	}
}
