package handler

import (
	"net/http"
	"strconv"

	"Gym-StrongCode/internal/service"

	"github.com/gin-gonic/gin"
)

type TrainerHandler struct {
	trainerService *service.TrainerService
}

func NewTrainerHandler(trainerService *service.TrainerService) *TrainerHandler {
	return &TrainerHandler{trainerService: trainerService}
}

func (h *TrainerHandler) List(c *gin.Context) {
	trainers, err := h.trainerService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trainers)
}

type createTrainerRequest struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

func (h *TrainerHandler) Create(c *gin.Context) {
	var req createTrainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trainer, err := h.trainerService.Create(req.Name, req.Bio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trainer)
}

type updateTrainerRequest struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

func (h *TrainerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req updateTrainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.trainerService.Update(id, req.Name, req.Bio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trainer updated"})
}

func (h *TrainerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.trainerService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trainer deleted"})
}
