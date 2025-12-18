package handler

import (
	"net/http"
	"strconv"

	"Gym_StrongCode/internal/service"

	"github.com/gin-gonic/gin"
)

type TrainerHandler struct {
	trainerService *service.TrainerService
}

func NewTrainerHandler(trainerService *service.TrainerService) *TrainerHandler {
	return &TrainerHandler{trainerService: trainerService}
}

// ListTrainers godoc
// @Summary      List all trainers
// @Description  Get all fitness trainers
// @Tags         trainers
// @Produce      json
// @Success      200  {array}   models.Trainer
// @Failure      500  {object}  map[string]string
// @Router       /trainers [get]
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

// CreateTrainer godoc
// @Summary      Create trainer
// @Description  Create a new trainer profile (admin only)
// @Tags         trainers
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        body  body      handler.createTrainerRequest  true  "Trainer data"
// @Success      201   {object}  models.Trainer
// @Failure      400   {object}  map[string]string
// @Router       /admin/trainers [post]
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

// UpdateTrainer godoc
// @Summary      Update trainer
// @Description  Update trainer profile (admin only)
// @Tags         trainers
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id    path      int                         true  "Trainer ID"
// @Param        body  body      handler.updateTrainerRequest true  "Updated trainer data"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /admin/trainers/{id} [put]
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

// DeleteTrainer godoc
// @Summary      Delete trainer
// @Description  Delete trainer profile (admin only)
// @Tags         trainers
// @Security     Bearer
// @Param        id   path      int  true  "Trainer ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /admin/trainers/{id} [delete]
func (h *TrainerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.trainerService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trainer deleted"})
}
