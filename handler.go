package main

import (
	"time"

	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
	ginmw "github.com/coinbase/x402/go/http/gin"
	evm "github.com/coinbase/x402/go/mechanisms/evm/exact/server"
	"github.com/ride4Low/contracts/env"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tripClient       *TripClient
	stripeWebhookKey string
	publisher        EventPublisher
}

func NewHandler(tripClient *TripClient, stripeWebhookKey string, publisher EventPublisher) *Handler {
	return &Handler{
		tripClient:       tripClient,
		stripeWebhookKey: stripeWebhookKey,
		publisher:        publisher,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	tripGroup := r.Group("/api/trip")
	tripGroup.POST("/preview", h.previewTrip)
	tripGroup.POST("/start", h.startTrip)

	tripGroup.POST("/pay", h.payTrip)

	r.Group("/webhook").POST("/stripe", h.handleWebhookStripe)
}

func (h *Handler) ApplyX402Middleware(r *gin.Engine) {
	evmAddress := env.GetString("EVM_PAYEE", "")
	facilitatorURL := env.GetString("FACILITATOR_ENDPOINT", "http://localhost:4022")
	network := x402.ParseNetwork("eip155:84532")

	routes := x402http.RoutesConfig{
		"POST /api/trip/pay": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  "exact",
					Price:   "$0.001",
					Network: network,
					PayTo:   evmAddress,
				},
			},
			Description: "Pay for a trip",
			MimeType:    "application/json",
		},
	}

	// Create HTTP facilitator client
	facilitatorClient := x402http.NewHTTPFacilitatorClient(&x402http.FacilitatorConfig{
		URL: facilitatorURL,
	})

	// Apply x402 payment middleware
	r.Use(ginmw.X402Payment(ginmw.Config{
		Routes:      routes,
		Facilitator: facilitatorClient,
		Schemes: []ginmw.SchemeConfig{
			{Network: network, Server: evm.NewExactEvmScheme()},
		},
		Timeout: 30 * time.Second,
	}))
}

func (h *Handler) ErrorResponse(c *gin.Context, httpCode int, code ErrorCodes, message string) {
	c.JSON(httpCode, APIResponse{
		Error: &APIError{
			Code:    string(code),
			Message: message,
		},
	})
}

func (h *Handler) DataResponse(c *gin.Context, httpCode int, data any) {
	c.JSON(httpCode, APIResponse{
		Data: data,
	})
}
