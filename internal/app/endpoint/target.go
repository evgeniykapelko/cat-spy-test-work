package endpoint

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"spy_cat/internal/app/model"
)

type TargetService interface {
	Create(target *model.Target) error
	FindAll() ([]model.Target, error)
	FindByID(id uint) (*model.Target, error)
	Update(target *model.Target) error
	Delete(target *model.Target) error
	CompleteTarget(targetID uint) error
}

type TargetEndpoint struct {
	s TargetService
}

func NewTargetEndpoint(s TargetService) *TargetEndpoint {
	return &TargetEndpoint{
		s: s,
	}
}

func (e *TargetEndpoint) CreateTarget(c echo.Context) error {
	var target model.Target
	if err := c.Bind(&target); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if err := e.s.Create(&target); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create target"})
	}
	return c.JSON(http.StatusOK, target)
}

func (e *TargetEndpoint) GetAllTargets(c echo.Context) error {
	targets, err := e.s.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve targets"})
	}
	return c.JSON(http.StatusOK, targets)
}

func (e *TargetEndpoint) GetTarget(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target ID"})
	}

	target, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
	}
	return c.JSON(http.StatusOK, target)
}

func (e *TargetEndpoint) UpdateTarget(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target ID"})
	}

	target, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
	}

	if err := c.Bind(target); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := e.s.Update(target); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update target"})
	}
	return c.JSON(http.StatusOK, target)
}

func (e *TargetEndpoint) DeleteTarget(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target ID"})
	}

	target, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
	}
	if err := e.s.Delete(target); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete target"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Target deleted successfully"})
}

func (e *TargetEndpoint) CompleteTarget(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target ID"})
	}

	if err := e.s.CompleteTarget(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to complete target"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Target completed successfully"})
}
