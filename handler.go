package main

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	tripClient *TripClient
}

func NewHandler(tripClient *TripClient) *Handler {
	return &Handler{
		tripClient: tripClient,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	tripGroup := r.Group("/api/trip")
	tripGroup.POST("/preview", h.previewTrip)
	tripGroup.POST("/start", h.startTrip)
}

func (h *Handler) ErrorResponse(c *gin.Context, httpCode int, code string, message string) {
	c.JSON(httpCode, APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func (h *Handler) DataResponse(c *gin.Context, httpCode int, data any) {
	c.JSON(httpCode, APIResponse{
		Data: data,
	})
}
