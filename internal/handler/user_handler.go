package handler

import (
	"net/http"
	"strconv"

	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

type updateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetCurrent godoc
// @Summary      Get current user
// @Description  Returns the currently authenticated user profile
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      500  {object}  map[string]string
// @Router       /me [get]
func (h *UserHandler) GetCurrent(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary      Update user profile
// @Description  Update name and email of current user
// @Tags         users
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        body  body      updateUserRequest  true  "Update data"
// @Success      200   {object}  models.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /me [put]
func (h *UserHandler) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req updateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userRepo.Update(userID, req.Name, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	user, _ := h.userRepo.GetByID(userID)
	c.JSON(http.StatusOK, user)
}

// ListUsers godoc
// @Summary      List all users
// @Description  Get all users (admin only)
// @Tags         users
// @Security     Bearer
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /admin/users [get]
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID (admin only)
// @Tags         users
// @Security     Bearer
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /admin/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.userRepo.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
