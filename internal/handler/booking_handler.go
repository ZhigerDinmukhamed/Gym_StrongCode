package handler

import (
	"Gym_StrongCode/internal/repository"
	"net/http"
	"strconv"

	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/service"

	"github.com/gin-gonic/gin"
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

// CreateBooking godoc
// @Summary      Book a class
// @Description  Create a booking for a fitness class
// @Tags         bookings
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        body  body      handler.createBookingRequest  true  "Booking data"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /bookings [post]
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

// ListUserBookings godoc
// @Summary      List user bookings
// @Description  Get all bookings for the current user
// @Tags         bookings
// @Security     Bearer
// @Produce      json
// @Success      200  {array}   models.Booking
// @Failure      500  {object}  map[string]string
// @Router       /bookings [get]
func (h *BookingHandler) ListUser(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	bookings, err := h.bookingService.ListUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// ListAllBookings godoc
// @Summary      List all bookings
// @Description  Get all bookings in the system (admin only)
// @Tags         bookings
// @Security     Bearer
// @Produce      json
// @Success      200  {array}   models.Booking
// @Failure      500  {object}  map[string]string
// @Router       /admin/bookings [get]
func (h *BookingHandler) ListAll(c *gin.Context) { // admin only
	bookings, err := h.bookingService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// CancelBooking godoc
// @Summary      Cancel booking
// @Description  Cancel a user's class booking
// @Tags         bookings
// @Security     Bearer
// @Param        id   path      int  true  "Booking ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /bookings/{id} [delete]
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
