package model

type Target struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Complete  bool   `json:"complete"`
	MissionID uint   `json:"mission_id"`
}
