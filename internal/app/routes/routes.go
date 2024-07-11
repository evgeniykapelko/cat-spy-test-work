package routes

import (
	"github.com/labstack/echo/v4"
	"spy_cat/internal/app/endpoint"
)

func RegisterRoutes(e *echo.Echo, catEndpoint *endpoint.CatEndpoint, missionEndpoint *endpoint.MissionEndpoint, targetEndpoint *endpoint.TargetEndpoint) {
	e.POST("/cats", catEndpoint.CreateCat)
	e.GET("/cats", catEndpoint.GetAllCats)
	e.GET("/cats/:id", catEndpoint.GetCat)
	e.PUT("/cats/:id/salary", catEndpoint.UpdateCatSalary)
	e.DELETE("/cats/:id", catEndpoint.DeleteCat)

	e.POST("/missions", missionEndpoint.CreateMission)
	e.GET("/missions", missionEndpoint.GetAllMissions)
	e.GET("/missions/:id", missionEndpoint.GetMission)
	e.PUT("/missions/:id", missionEndpoint.UpdateMission)
	e.DELETE("/missions/:id", missionEndpoint.DeleteMission)
	e.PUT("/missions/:missionID/cats/:catID", missionEndpoint.AssignCatToMission)
	e.PUT("/missions/:id/complete", missionEndpoint.CompleteMission)

	e.PATCH("/missions/:missionID/targets/:targetID/notes", missionEndpoint.UpdateTargetNotes)
	e.POST("/missions/:id/targets", missionEndpoint.AddTargetToMission)
	e.DELETE("/missions/:missionId/targets/:targetId", missionEndpoint.DeleteTargetFromMission)

	e.POST("/targets", targetEndpoint.CreateTarget)
	e.GET("/targets", targetEndpoint.GetAllTargets)
	e.GET("/targets/:id", targetEndpoint.GetTarget)
	e.PUT("/targets/:id", targetEndpoint.UpdateTarget)
	e.DELETE("/targets/:id", targetEndpoint.DeleteTarget)
	e.PUT("/targets/:id/complete", targetEndpoint.CompleteTarget)

}
