package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func enableCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, PAYMENT-SIGNATURE, Price, TripID")
	c.Header("Access-Control-Expose-Headers", "PAYMENT-REQUIRED")

	// allow preflight requests from the browser API
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}
