package lib

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func DefaultServerEngine() *chi.Mux {
	return chi.NewRouter()
}

func RunServer[S, C, R any](config ServerConfig[S, C, R], port int, debug bool, logger Logger) Error {
	var err Error

	var engine *chi.Mux
	if config.InitServerEngine != nil {
		engine, err = config.InitServerEngine(config.Storage)
		if err != nil {
			return err.Tap()
		}
	} else {
		engine = DefaultServerEngine()
	}

	routesManager := NewRouterManager[S, C, R](
		engine,
		config.Storage,
		config.MakeContext,
		config.CloseContext,
		debug,
		logger,
	)
	config.RegistrationRoutes(routesManager)

	fmt.Printf("Starting http_server at port %d...", port)
	err_ := http.ListenAndServe(fmt.Sprintf(":%d", port), engine)
	if err_ != nil {
		return Err(fmt.Errorf("error starting http server: %v", err))
	}

	return nil
}
