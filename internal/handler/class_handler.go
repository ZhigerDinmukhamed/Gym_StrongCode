package handler

import (
	"net/http"
	"strconv"

	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/service"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	classService *service.ClassService
}

func NewClassHandler(classService *service.ClassService) *ClassHandler {
	return &ClassHandler{classService: classService}
}

// ListClasses godoc
// @Summary      List all classes
// @Description  Get all fitness classes
// @Tags         classes
// @Produce      json
// @Success      200  {array}   models.Class
// @Failure      500  {object}  map[string]string
// @Router       /classes [get]
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

// CreateClass godoc
// @Summary      Create class
// @Description  Create a new fitness class (admin only)
// @Tags         classes
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        body  body      handler.createClassRequest  true  "Class data"
// @Success      201   {object}  models.Class
// @Failure      400   {object}  map[string]string
// @Router       /admin/classes [post]
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

// UpdateClass godoc
// @Summary      Update class
// @Description  Update fitness class details (admin only)
// @Tags         classes
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id    path      int                         true  "Class ID"
// @Param        body  body      handler.createClassRequest  true  "Updated class data"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /admin/classes/{id} [put]
func (h *ClassHandler) Update(c *gin.Context) {
	// аналогично Create, с парсингом id
	// ... (реализация похожа)
}

// DeleteClass godoc
// @Summary      Delete class
// @Description  Delete fitness class (admin only)
// @Tags         classes
// @Security     Bearer
// @Param        id   path      int  true  "Class ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /admin/classes/{id} [delete]
func (h *ClassHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.classService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "class deleted"})
}
