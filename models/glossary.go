package models

type Glossary struct {
	GlossaryEn  string `json:"glossaryEn"`
	GlossaryAr  string `json:"glossaryAr"`
	GlossaryURL string `json:"glossaryUrl"`
	License     string `json:"license,omitempty"`
}
