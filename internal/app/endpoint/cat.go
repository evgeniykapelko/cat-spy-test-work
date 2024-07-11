package endpoint

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"spy_cat/internal/app/model"
	"spy_cat/internal/app/service"
)

type CatService interface {
	Create(cat *model.Cat) error
	FindAll() ([]model.Cat, error)
	FindByID(id uint) (*model.Cat, error)
	Update(cat *model.Cat) error
	Delete(cat *model.Cat) error
}

type CatEndpoint struct {
	s              CatService
	breedValidator *service.BreedValidator
}

func NewCatEndpoint(s CatService) *CatEndpoint {
	return &CatEndpoint{
		s:              s,
		breedValidator: service.NewBreedValidator(),
	}
}

func (e *CatEndpoint) CreateCat(c echo.Context) error {
	var cat model.Cat
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := e.breedValidator.Validate(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := e.s.Create(&cat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create cat"})
	}
	return c.JSON(http.StatusOK, cat)
}

func (e *CatEndpoint) GetAllCats(c echo.Context) error {
	cats, err := e.s.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve cats"})
	}
	return c.JSON(http.StatusOK, cats)
}

func (e *CatEndpoint) GetCat(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cat ID"})
	}

	cat, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
	}
	return c.JSON(http.StatusOK, cat)
}

func (e *CatEndpoint) UpdateCatSalary(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cat ID"})
	}

	cat, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
	}

	var input struct {
		Salary float64 `json:"salary"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	cat.Salary = input.Salary

	if err := e.s.Update(cat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update cat"})
	}
	return c.JSON(http.StatusOK, cat)
}

func (e *CatEndpoint) DeleteCat(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cat ID"})
	}

	cat, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
	}
	if err := e.s.Delete(cat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete cat"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Cat deleted successfully"})
}
