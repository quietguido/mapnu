package rest

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/quietguido/mapnu/internal/services"
	"github.com/quietguido/mapnu/pkg/middleware"
)

type restH struct {
	lg       *zap.Logger
	services *services.Service
}

func initRest(lg *zap.Logger, services *services.Service) *restH {
	return &restH{
		lg:       lg,
		services: services,
	}
}

func GetHandler(
	lg *zap.Logger,
	services *services.Service,
) http.Handler {
	restH := initRest(lg, services)
	router := http.NewServeMux()

	lgMiddleware := middleware.NewLogging(lg)
	middlewareStack := middleware.CreateStack(lgMiddleware.Logging)

	router.HandleFunc("POST /event", restH.CreateEventHandler)
	router.HandleFunc("GET /event/{id}", restH.GetEventByIdHandler)

	return middlewareStack(router)
}
