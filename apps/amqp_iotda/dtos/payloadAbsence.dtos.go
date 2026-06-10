package dtos

type PayloadAbsence struct {
	Nisn       string `json:"nisn"`
	IsMatching bool   `json:"is_matching"`
}
