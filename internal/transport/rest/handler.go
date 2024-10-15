package rest

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/quietguido/mapnu/pkg/middleware"
)

func GetHandler(
	lg *zap.Logger,
) http.Handler {
	router := http.NewServeMux()

	lgMiddleware := middleware.NewLogging(lg)
	middlewareStack := middleware.CreateStack(lgMiddleware.Logging)

	router.HandleFunc("login", nil)

	return middlewareStack(router)
}
