package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	appServices "github.com/ZaneH/keep-talking/internal/application/services"
	domainServices "github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/infrastructure/adapters"
	"github.com/ZaneH/keep-talking/internal/infrastructure/rest"
)

func main() {
	actorSystem := actors.NewActorSystem()
	actorSystemAdapter := adapters.NewActorSystemAdapter(actorSystem)
	bombFactory := &domainServices.BombFactoryImpl{}
	bombService := appServices.NewBombService(actorSystemAdapter, bombFactory)
	gameService := appServices.NewGameService(actorSystem, bombService)

	restAdapter := rest.NewGameServiceAdapter(gameService)

	router := restAdapter.Router()
	router.Use(rest.LoggingMiddleware)
	router.Use(rest.CORSMiddleware)
	router.Use(rest.ContentTypeMiddleware)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("REST server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
}
