package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"Gym-StrongCode/internal/models"
	"Gym-StrongCode/internal/service"
)

type ClassHandler struct {
	classService *service.ClassService
}

func NewClassHandler(classService *service.ClassService) *ClassHandler {
	return &ClassHandler{classService: classService}
}

func (h *ClassHandler) List(c *gin.Context) {
	classes, err := h.classService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

type createClassRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	TrainerID   int    `json:"trainer_id"`
	GymID       int    `json:"gym_id" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	DurationMin int    `json:"duration_min" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required"`
}

func (h *ClassHandler) Create(c *gin.Context) {
	var req createClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class := &models.Class{
		Title:       req.Title,
		Description: req.Description,
		TrainerID:   req.TrainerID,
		GymID:       req.GymID,
		StartTime:   req.StartTime,
		DurationMin: req.DurationMin,
		Capacity:    req.Capacity,
	}

	created, err := h.classService.Create(class)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *ClassHandler) Update(c *gin.Context) {
	// аналогично Create, с парсингом id
	// ... (реализация похожа)
}

func (h *ClassHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.classService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "class deleted"})
}