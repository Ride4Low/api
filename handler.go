package main

import (
	"context"
	"errors"
	"log"
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

	dynamicPrice := func(ctx context.Context, hc x402http.HTTPRequestContext) (x402.Price, error) {

		adapter, ok := hc.Adapter.(*ginmw.GinAdapter)
		if !ok {
			return nil, errors.New("not a gin adapter")
		}

		var x402Price x402.Price

		tripID := adapter.GetHeader("TripID")
		log.Println(tripID)
		// find trip in trip service to get price
		// or get jwt payload and validate

		// for simple usage, I'm just gonna extract Price from URL
		price := adapter.GetHeader("Price")
		log.Println(price)

		if price == "" || price == "0" {
			return nil, errors.New("price is empty or zero")
		}

		x402Price = price

		return x402Price, nil
	}

	routes := x402http.RoutesConfig{
		"POST /api/trip/pay": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme: "exact",
					// Price:   "$0.001",
					Price:   x402http.DynamicPriceFunc(dynamicPrice),
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
