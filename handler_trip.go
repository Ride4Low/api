package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ride4Low/contracts/proto/trip"
)

func (h *Handler) previewTrip(c *gin.Context) {
	var req previewTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.tripClient.tripClient.PreviewTrip(c.Request.Context(), &trip.PreviewTripRequest{
		PickupLocation: &trip.Coordinate{
			Latitude:  req.Pickup.Latitude,
			Longitude: req.Pickup.Longitude,
		},
		DropoffLocation: &trip.Coordinate{
			Latitude:  req.Destination.Latitude,
			Longitude: req.Destination.Longitude,
		},
	})

	var apiResponse APIResponse

	if err != nil {
		apiResponse.Error = &APIError{
			Code:    string(PreviewTripError),
			Message: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, apiResponse)
		return
	}

	h.DataResponse(c, http.StatusOK, response)
}

func (h *Handler) startTrip(c *gin.Context) {

}
