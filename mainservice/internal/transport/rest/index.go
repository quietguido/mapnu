package rest

import (
	"github.com/quietguido/mapnu/mainservice/internal/services/oauth"
	"net/http"

	"go.uber.org/zap"

	"github.com/quietguido/mapnu/mainservice/internal/services"
	"github.com/quietguido/mapnu/mainservice/pkg/middleware"
)

type restH struct {
	lg       *zap.Logger
	services *services.Service
	oauthH   *OAuthHandler
}

func initRest(lg *zap.Logger, services *services.Service) *restH {
	oauthService := oauth.NewOAuthService(lg)
	oauthHandler := NewOAuthHandler(oauthService)

	return &restH{
		lg:       lg,
		services: services,
		oauthH:   oauthHandler,
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

	//oauth
	router.HandleFunc("/api/user/profile", restH.oauthH.GetUserProfile)
	router.HandleFunc("POST /auth/token/exchange", restH.oauthH.HandleTokenExchange)

	return middlewareStack(router)
}
