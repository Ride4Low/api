package main

import (
	"github.com/ride4Low/contracts/proto/trip"
	"github.com/ride4Low/contracts/types"
)

type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type previewTripRequest struct {
	UserID  string           `json:"userID" binding:"required"`
	Pickup  types.Coordinate `json:"pickup" binding:"required"`
	Dropoff types.Coordinate `json:"dropoff" binding:"required"`
}

func (p *previewTripRequest) toProto() *trip.PreviewTripRequest {
	return &trip.PreviewTripRequest{
		UserID: p.UserID,
		PickupLocation: &trip.Coordinate{
			Latitude:  p.Pickup.Latitude,
			Longitude: p.Pickup.Longitude,
		},
		DropoffLocation: &trip.Coordinate{
			Latitude:  p.Dropoff.Latitude,
			Longitude: p.Dropoff.Longitude,
		},
	}
}
