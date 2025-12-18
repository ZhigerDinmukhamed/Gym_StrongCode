package handler

import (
	"net/http"
	"strconv"

	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/service"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

// ListMemberships godoc
// @Summary      List all memberships
// @Description  Get all available membership plans
// @Tags         memberships
// @Produce      json
// @Success      200  {array}   models.Membership
// @Failure      500  {object}  map[string]string
// @Router       /memberships [get]
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

// BuyMembership godoc
// @Summary      Buy membership
// @Description  Purchase a membership plan
// @Tags         memberships
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        body  body      handler.buyMembershipRequest  true  "Purchase data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Router       /memberships/buy [post]
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

// CreateMembership godoc
// @Summary      Create membership plan
// @Description  Create a new membership plan (admin only)
// @Tags         memberships
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        body  body      handler.createMembershipRequest  true  "Membership data"
// @Success      201   {object}  models.Membership
// @Failure      400   {object}  map[string]string
// @Router       /admin/memberships [post]
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

// UpdateMembership godoc
// @Summary      Update membership plan
// @Description  Update membership plan details (admin only)
// @Tags         memberships
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id    path      int                             true  "Membership ID"
// @Param        body  body      handler.createMembershipRequest true  "Updated membership data"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /admin/memberships/{id} [put]
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

// DeleteMembership godoc
// @Summary      Delete membership plan
// @Description  Delete membership plan (admin only)
// @Tags         memberships
// @Security     Bearer
// @Param        id   path      int  true  "Membership ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /admin/memberships/{id} [delete]
func (h *MembershipHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.membershipService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "membership deleted"})
}
