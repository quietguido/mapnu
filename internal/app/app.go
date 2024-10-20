package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/quietguido/mapnu/internal/database/psql"
	"github.com/quietguido/mapnu/internal/repo"
	"github.com/quietguido/mapnu/internal/services"
	"github.com/quietguido/mapnu/internal/transport/rest"
	"github.com/quietguido/mapnu/pkg/httpserver"
)

func Execute() {
	err := godotenv.Load("config.env")
	if err != nil {
		panic(err.Error() + " failed to load config.env")
	}

	lg, err := zap.NewProduction()
	if err != nil {
		panic("lg creation error")
	}

	dbcon, err := psql.New(psql.Config{
		Addr:     os.Getenv("POSTGRES_HOST"), //change for local and docker
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DB:       os.Getenv("POSTGRES_DB"),
	})
	if err != nil {
		panic(err)
	}

	repos := repo.InitRepositories(lg, dbcon)
	services := services.InitServices(lg, repos)
	restHandler := rest.GetHandler(lg, services)
	server := httpserver.New(":8080", restHandler)

	oschan := make(chan os.Signal, 1)
	signal.Notify(oschan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// start server
	lg.Log(zapcore.InfoLevel, "start server")
	server.Start()
	lg.Log(zapcore.InfoLevel, "started server")

	var exitcode int

	select {
	case <-oschan:
	case <-server.Wait():
		exitcode = 1
	}

	// Gracefulshutdown

	if err = server.Shutdown(20 * time.Second); err != nil {
		exitcode = 1
	}

	lg.Log(zapcore.InfoLevel, "exit")
	os.Exit(exitcode)
}
