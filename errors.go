package main

// error codes
type ErrorCodes string

const (
	// Request Level Errors
	InvalidRequestError ErrorCodes = "INVALID_REQUEST_ERROR"

	// Service Level Errors
	PreviewTripError ErrorCodes = "PREVIEW_TRIP_ERROR"
	StartTripError   ErrorCodes = "START_TRIP_ERROR"
	CreateTripError  ErrorCodes = "CREATE_TRIP_ERROR"
)
