package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"Gym-StrongCode/internal/middleware"
	"Gym-StrongCode/internal/service"
)

type GymHandler struct {
	gymService *service.GymService
}

func NewGymHandler(gymService *service.GymService) *GymHandler {
	return &GymHandler{gymService: gymService}
}

func (h *GymHandler) List(c *gin.Context) {
	gyms, err := h.gymService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gyms)
}

type createGymRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func (h *GymHandler) Create(c *gin.Context) {
	var req createGymRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gym, err := h.gymService.Create(req.Name, req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gym)
}

type updateGymRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
}

func (h *GymHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req updateGymRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.gymService.Update(id, req.Name, req.Address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "gym updated"})
}

func (h *GymHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.gymService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "gym deleted"})
}