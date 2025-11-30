package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) previewTrip(c *gin.Context) {
	var req previewTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, InvalidRequestError, err.Error())
		return
	}

	response, err := h.tripClient.tripClient.PreviewTrip(c.Request.Context(), req.toProto())
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, PreviewTripError, err.Error())
		return
	}

	h.DataResponse(c, http.StatusOK, response)
}

func (h *Handler) startTrip(c *gin.Context) {

}
