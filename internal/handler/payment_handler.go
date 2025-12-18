package handler

import (
	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
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

// CreatePayment godoc
// @Summary      Create payment
// @Description  Create a standalone payment
// @Tags         payments
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        body  body      handler.createPaymentRequest  true  "Payment data"
// @Success      201   {object}  models.Payment
// @Failure      400   {object}  map[string]string
// @Router       /payments [post]
func (h *PaymentHandler) CreateStandalone(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentService.Create(userID, req.AmountCents, req.Method, "", "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// ListAllPayments godoc
// @Summary      List all payments
// @Description  Get all payments in the system (admin only)
// @Tags         payments
// @Security     Bearer
// @Produce      json
// @Success      200  {array}   models.Payment
// @Failure      500  {object}  map[string]string
// @Router       /admin/payments [get]
func (h *PaymentHandler) ListAll(c *gin.Context) {
	payments, err := h.paymentService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
