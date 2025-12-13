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
	"github.com/ride4Low/contracts/pkg/otel"
	"github.com/ride4Low/contracts/pkg/rabbitmq"
)

var (
	httpAddr       = env.GetString("HTTP_ADDR", ":8081")
	tripServiceURL = env.GetString("TRIP_SERVICE_ADDR", "localhost:9093")
	rabbitMQURI    = env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	jaegerEndpoint = env.GetString("JAEGER_ENDPOINT", "http://localhost:4317")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup OpenTelemetry
	otelCfg := otel.DefaultConfig("api")
	otelCfg.JaegerEndpoint = jaegerEndpoint
	otelProvider, err := otel.Setup(ctx, otelCfg)
	if err != nil {
		log.Fatalf("Failed to setup OpenTelemetry: %v", err)
	}
	defer func() {
		if err := otelProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down OpenTelemetry provider: %v", err)
		}
	}()

	r := gin.Default()

	r.Use(otel.GinMiddleware("api"))
	r.Use(enableCORS)

	tripClient, err := NewTripClient(tripServiceURL)
	if err != nil {
		log.Fatalf("Failed to create trip client: %v", err)
	}
	defer tripClient.Close()

	rmq, err := rabbitmq.NewRabbitMQ(rabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ: %v", err)
	}
	defer rmq.Close()

	rmqPublisher := rabbitmq.NewPublisher(rmq)
	eventPublisher := NewAmqpPublisher(rmqPublisher)

	h := NewHandler(tripClient, env.GetString("STRIPE_WEBHOOK_KEY", ""), eventPublisher)
	h.ApplyX402Middleware(r)
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
