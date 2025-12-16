package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"Gym-StrongCode/internal/middleware"
	"Gym-StrongCode/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type createPaymentRequest struct {
	AmountCents int    `json:"amount_cents" binding:"required,gt=0"`
	Method      string `json:"method" binding:"required"`
}

func (h *PaymentHandler) CreateStandalone(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.CreateStandalone(userID, req.AmountCents, req.Method)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) ListAll(c *gin.Context) { // admin only
	payments, err := h.paymentService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}