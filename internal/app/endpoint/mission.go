package endpoint

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"spy_cat/internal/app/model"
)

type MissionService interface {
	Create(mission *model.Mission) error
	FindAll() ([]model.Mission, error)
	FindByID(id uint) (*model.Mission, error)
	Update(mission *model.Mission) error
	Delete(mission *model.Mission) error
	AssignCatToMission(missionID uint, catID uint) error
	CompleteMission(missionID uint) error
	UpdateTargetNotes(missionId, targetId uint, notes string) error
	FindTargetByID(id uint) (*model.Target, error)
	DeleteTarget(target *model.Target) error
	AddTarget(mission *model.Mission, target *model.Target) error
}

type MissionEndpoint struct {
	s MissionService
}

func NewMissionEndpoint(s MissionService) *MissionEndpoint {
	return &MissionEndpoint{
		s: s,
	}
}

func (e *MissionEndpoint) CreateMission(c echo.Context) error {
	var mission model.Mission
	if err := c.Bind(&mission); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	if err := e.s.Create(&mission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create mission"})
	}
	return c.JSON(http.StatusOK, mission)
}

func (e *MissionEndpoint) GetAllMissions(c echo.Context) error {
	missions, err := e.s.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve missions"})
	}
	return c.JSON(http.StatusOK, missions)
}

func (e *MissionEndpoint) GetMission(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	mission, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
	}
	return c.JSON(http.StatusOK, mission)
}

func (e *MissionEndpoint) UpdateMission(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	mission, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
	}

	if err := c.Bind(mission); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := e.s.Update(mission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update mission"})
	}
	return c.JSON(http.StatusOK, mission)
}

func (e *MissionEndpoint) DeleteMission(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	mission, err := e.s.FindByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
	}

	if mission.CatID != 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete mission assigned to a cat"})
	}

	if err := e.s.Delete(mission); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete mission"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Mission deleted successfully"})
}

func (e *MissionEndpoint) AssignCatToMission(c echo.Context) error {
	missionIDStr := c.Param("missionID")
	missionID, err := strconv.ParseUint(missionIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	catIDStr := c.Param("catID")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cat ID"})
	}

	if err := e.s.AssignCatToMission(uint(missionID), uint(catID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign cat to mission"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Cat assigned to mission successfully"})
}

func (e *MissionEndpoint) CompleteMission(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	if err := e.s.CompleteMission(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to complete mission"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Mission completed successfully"})
}

func (e *MissionEndpoint) UpdateTargetNotes(c echo.Context) error {
	missionId, err := strconv.ParseUint(c.Param("missionID"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	targetId, err := strconv.ParseUint(c.Param("targetID"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target ID"})
	}

	mission, err := e.s.FindByID(uint(missionId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
	}
	if mission.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Mission is complete, notes cannot be updated"})
	}

	target, err := e.s.FindTargetByID(uint(targetId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
	}
	if target.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Target is complete, notes cannot be updated"})
	}

	var input struct {
		Notes string `json:"notes"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := e.s.UpdateTargetNotes(uint(missionId), uint(targetId), input.Notes); err != nil {
		if err.Error() == "mission is complete, notes cannot be updated" || err.Error() == "target is complete, notes cannot be updated" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update notes"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Notes updated successfully"})
}

func (e *MissionEndpoint) AddTargetToMission(c echo.Context) error {
	missionIdStr := c.Param("id")
	missionId, err := strconv.ParseUint(missionIdStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	mission, err := e.s.FindByID(uint(missionId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
	}
	if mission.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot add target to a completed mission"})
	}

	var target model.Target
	if err := c.Bind(&target); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	target.MissionID = uint(missionId)

	if err := e.s.AddTarget(mission, &target); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add target"})
	}

	return c.JSON(http.StatusCreated, target)
}

func (e *MissionEndpoint) DeleteTargetFromMission(c echo.Context) error {
	targetIdStr := c.Param("targetId")
	targetId, err := strconv.ParseUint(targetIdStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid target ID"})
	}

	target, err := e.s.FindTargetByID(uint(targetId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Target not found"})
	}
	if target.Complete {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete a completed target"})
	}

	if err := e.s.DeleteTarget(target); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete target"})
	}

	return c.JSON(http.StatusNoContent, nil)
}
