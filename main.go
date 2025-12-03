package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ride4Low/contracts/env"
)

var (
	httpAddr       = env.GetString("HTTP_ADDR", ":8081")
	tripServiceURL = env.GetString("TRIP_SERVICE_ADDR", "localhost:9093")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	r := gin.Default()

	r.Use(enableCORS)

	tripClient, err := NewTripClient(tripServiceURL)
	if err != nil {
		log.Fatalf("Failed to create trip client: %v", err)
	}
	defer tripClient.Close()

	h := NewHandler(tripClient)
	h.RegisterRoutes(r)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully: %v", err)
			server.Close()
		}
	}
}
