package main

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	tripGroup := r.Group("/api/trip")
	tripGroup.POST("/preview", h.previewTrip)
	tripGroup.POST("/start", h.startTrip)
}

func (h *Handler) previewTrip(c *gin.Context) {

}

func (h *Handler) startTrip(c *gin.Context) {

}
