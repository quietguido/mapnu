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

	//user
	router.HandleFunc("POST /user", restH.CreateUserHandler)
	router.HandleFunc("GET /user/{id}", restH.GetUserByIdHandler)

	//event
	router.HandleFunc("POST /event", restH.CreateEventHandler)
	router.HandleFunc("GET /event/{id}", restH.GetEventByIdHandler)
	router.HandleFunc("GET /map", restH.GetMapForQuadrantHandler)

	//booking
	router.HandleFunc("POST /booking", restH.CreateBookingHandler)
	router.HandleFunc("GET /booking/{id}", restH.GetBookingByIdHandler)
	router.HandleFunc("GET /booking", restH.GetBookingsForUserHandler)
	router.HandleFunc("POST /booking/status", restH.ChangeBookingStatusHandler)
	router.HandleFunc("GET /booking/organizer", restH.GetBookingApplicationsForOrganizer)

	return middlewareStack(router)
}
