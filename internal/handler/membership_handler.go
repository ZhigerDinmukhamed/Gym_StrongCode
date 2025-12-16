package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"Gym-StrongCode/internal/middleware"
	"Gym-StrongCode/internal/service"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

func (h *MembershipHandler) List(c *gin.Context) {
	memberships, err := h.membershipService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memberships)
}

type buyMembershipRequest struct {
	MembershipID int    `json:"membership_id" binding:"required"`
	Method       string `json:"method" binding:"required"`
}

func (h *MembershipHandler) Buy(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req buyMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.membershipService.Buy(userID, req.MembershipID, req.Method)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Admin CRUD
type createMembershipRequest struct {
	Name         string `json:"name" binding:"required"`
	DurationDays int    `json:"duration_days" binding:"required"`
	PriceCents   int    `json:"price_cents" binding:"required"`
}

func (h *MembershipHandler) Create(c *gin.Context) {
	var req createMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.membershipService.Create(req.Name, req.DurationDays, req.PriceCents)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, m)
}

func (h *MembershipHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req createMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.membershipService.Update(id, req.Name, req.DurationDays, req.PriceCents); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "membership updated"})
}

func (h *MembershipHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.membershipService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "membership deleted"})
}