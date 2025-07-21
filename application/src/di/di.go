package di

import (
	"challengephp/lib"
	"challengephp/src/services/repository"
	"challengephp/src/services/response_service"
	"challengephp/src/storage"
	"context"
	"net/http"
)

type Container struct {
	s   *storage.Storage
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}

func MakeContainer(s *storage.Storage, w http.ResponseWriter, r *http.Request) *Container {
	return &Container{
		s:   s,
		w:   w,
		r:   r,
		ctx: context.Background(),
	}
}

func CloseContainer(_ *Container) {
}

func (it *Container) ResponseFactory() response_service.ResponseFactory {
	return response_service.NewResponseFactory()
}

func (it *Container) Debug() bool {
	return it.s.Debug
}

func (it *Container) Logger() lib.Logger {
	return it.s.Logger
}

func (it *Container) ResponseWriter() http.ResponseWriter {
	return it.w
}

func (it *Container) Request() *http.Request {
	return it.r
}

func (it *Container) Context() context.Context {
	return it.ctx
}

func (it *Container) Repository() repository.Repository {
	return repository.NewRepository(it.s.Pool, it.Context())
}
