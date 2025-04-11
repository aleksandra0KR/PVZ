package main

import (
	"context"
	"errors"
	"final/internal/controller"
	"final/internal/database"
	"final/internal/logger"
	auth "final/internal/midleware"
	"final/internal/repository"
	"final/internal/usecase"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Error("error loading .env file")
	}

	port := os.Getenv("HTTP_PORT")
	db := database.InitializeDBPostgres(3, 10)
	logger.InitLogger()
	err := auth.InitAuthFromConfig()
	if err != nil {
		log.Error(err)
	}

	repository := repository.NewRepository(db.GetDB())
	usecase := usecase.NewUseCase(repository)
	handlers := controller.NewHandler(*usecase)
	router := handlers.Handle()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Infof("server is running on port %s", port)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Info("timeout of 5 seconds.")
	}
	log.Info("Server exiting")
}
