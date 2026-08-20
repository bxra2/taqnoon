package models

type Publisher struct {
	ID           int64      `json:"id"`
	PublisherEn  string     `json:"publisherEn"`
	PublisherAr  string     `json:"publisherAr"`
	PublisherURL string     `json:"publisherUrl"`
	Glossaries   []Glossary `json:"glossaries"`
}
