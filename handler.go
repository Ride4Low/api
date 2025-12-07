package main

import (
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

	r.Group("/webhook").POST("/stripe", h.handleWebhookStripe)
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
