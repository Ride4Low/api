package main

import (
	"log"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/ride4Low/contracts/events"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
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
	var req startTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, InvalidRequestError, err.Error())
		return
	}

	trip, err := h.tripClient.tripClient.CreateTrip(c.Request.Context(), req.toProto())
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, CreateTripError, err.Error())
		return
	}

	h.DataResponse(c, http.StatusOK, trip)
}

func (h *Handler) handleWebhookStripe(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, InvalidRequestError, err.Error())
		return
	}

	event, err := webhook.ConstructEventWithOptions(
		body,
		c.Request.Header.Get("Stripe-Signature"),
		h.stripeWebhookKey,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)

	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		h.ErrorResponse(c, http.StatusBadRequest, InvalidRequestError, err.Error())
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		if err := sonic.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("Error unmarshalling checkout session: %v", err)
			h.ErrorResponse(c, http.StatusBadRequest, InvalidRequestError, err.Error())
			return
		}

		log.Printf("Checkout session: %#v\n", session)

		if session.PaymentStatus != "paid" {
			log.Printf("Payment not paid: %v", session.PaymentStatus)
			h.ErrorResponse(c, http.StatusBadRequest, InvalidRequestError, "payment not paid")
			return
		}

		payload := events.PaymentStatusUpdateData{
			TripID:   session.Metadata["trip_id"],
			UserID:   session.Metadata["user_id"],
			DriverID: session.Metadata["driver_id"],
		}

		err := h.publisher.PublishPaymentSuccess(c.Request.Context(), &payload)
		if err != nil {
			log.Printf("Error publishing payment success event: %v", err)
			h.ErrorResponse(c, http.StatusInternalServerError, InvalidRequestError, err.Error())
			return
		}

	}
}

func (h *Handler) payTrip(c *gin.Context) {
	tripID := c.GetHeader("TripID")
	if tripID == "" {
		h.ErrorResponse(c, http.StatusBadRequest, InvalidRequestError, "trip ID is required")
		return
	}
	log.Println("TripID", tripID)

	payload := events.PaymentStatusUpdateData{
		TripID: tripID,
	}

	err := h.publisher.PublishPaymentSuccess(c.Request.Context(), &payload)
	if err != nil {
		log.Printf("Error publishing payment success event: %v", err)
		h.ErrorResponse(c, http.StatusInternalServerError, InvalidRequestError, err.Error())
		return
	}
	h.DataResponse(c, http.StatusOK, gin.H{
		"tripID": tripID,
	})

}
