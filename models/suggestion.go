package models

type Suggestion struct {
	ID      int64  `json:"id"`
	English string `json:"english"`
	Arabic  string `json:"arabic"`
}
