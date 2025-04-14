package main

import (
	"context"
	"errors"
	"final/internal/controller"
	"final/internal/database"
	"final/internal/grpcPVZ"
	"final/internal/logger"
	"final/internal/metrics"
	auth "final/internal/midleware"
	"final/internal/repository"
	"final/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	logger.InitLogger()

	if err := godotenv.Load(); err != nil {
		log.Error("error loading .env file")
	}
	metrics.Init()

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	db := database.InitializeDBPostgres(3, 10)
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

	go func() {
		grpcPort := os.Getenv("GRPS_PORT")
		if grpcPort == "" {
			grpcPort = "3000"
		}
		grpcListener, errGrpc := net.Listen("tcp", ":"+grpcPort)
		if errGrpc != nil {
			log.Fatalf("failed to listen on port %s: %v", grpcPort, errGrpc)
		}

		grpcServer := grpc.NewServer()
		pvzService := grpcPVZ.NewPVZService(repository)
		grpcPVZ.RegisterPVZServiceServer(grpcServer, pvzService)
		reflection.Register(grpcServer)

		log.Infof("gRPC server is running on port %s", grpcPort)
		defer grpcServer.GracefulStop()
		if errGrpc = grpcServer.Serve(grpcListener); errGrpc != nil {
			log.Fatalf("failed to serve: %v", errGrpc)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		portMetrics := os.Getenv("METRICS_PORT")
		if portMetrics == "" {
			portMetrics = "9000"
		}

		errMetrics := http.ListenAndServe(":"+portMetrics, nil)
		if errMetrics != nil {
			log.Fatalf("Prometheus server failed: %v", err)
		}
		log.Infof("Prometheus server is running on port %s", portMetrics)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Info("timeout of 5 seconds.")
	}
	log.Info("Server exiting")
}
