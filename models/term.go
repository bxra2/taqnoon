package models

type Term struct {
	ID          int64  `json:"id"`
	English     string `json:"english"`
	Arabic      string `json:"arabic"`
	French      string `json:"french,omitempty"`
	German      string `json:"german,omitempty"`
	TURL        string `json:"tURL,omitempty"`
	Description string `json:"description,omitempty"`
	GlossaryEn  string `json:"glossaryEn"`
	GlossaryAr  string `json:"glossaryAr"`
	GlossaryURL string `json:"glossaryUrl"`
	PublisherID int64  `json:"publisherId"`
	PublisherAr  string `json:"publisherAr"`
	PublisherEn  string `json:"publisherEn"`
	PublisherURL string `json:"publisherUrl"`
}
