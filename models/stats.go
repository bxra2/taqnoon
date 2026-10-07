package models

type Stats struct {
	TotalTerms int64 `json:"totalTerms"`
	English    int64 `json:"english"`
	Arabic     int64 `json:"arabic"`
	French     int64 `json:"french"`
	German     int64 `json:"german"`
	Glossaries int64 `json:"glossaries"`
}
