package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"Gym-StrongCode/internal/middleware"
	"Gym-StrongCode/internal/service"
)

type BookingHandler struct {
	bookingService *service.BookingService
	userRepo       *repository.UserRepository // для получения email пользователя
}

func NewBookingHandler(bookingService *service.BookingService, userRepo *repository.UserRepository) *BookingHandler {
	return &BookingHandler{bookingService: bookingService, userRepo: userRepo}
}

type createBookingRequest struct {
	ClassID int `json:"class_id" binding:"required"`
}

func (h *BookingHandler) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req createBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем email пользователя для уведомления
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	if err := h.bookingService.Create(userID, req.ClassID, user.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "booking created"})
}

func (h *BookingHandler) ListUser(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	bookings, err := h.bookingService.ListUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) ListAll(c *gin.Context) { // admin only
	bookings, err := h.bookingService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.bookingService.Cancel(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}