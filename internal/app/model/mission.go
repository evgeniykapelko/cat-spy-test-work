package model

type Mission struct {
	ID       uint     `json:"id"`
	CatID    uint     `json:"cat_id"`
	Targets  []Target `json:"targets"`
	Complete bool     `json:"complete"`
}
