package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Glossary struct {
	GlossaryEn  string `bson:"glossaryEn"  json:"glossaryEn"`
	GlossaryAr  string `bson:"glossaryAr"  json:"glossaryAr"`
	GlossaryURL string `bson:"glossaryUrl" json:"glossaryUrl"`
}

type Publisher struct {
	ID           bson.ObjectID `bson:"_id,omitempty"  json:"id"`
	PublisherEn  string        `bson:"publisherEn"    json:"publisherEn"`
	PublisherAr  string        `bson:"publisherAr"    json:"publisherAr"`
	PublisherURL string        `bson:"publisherUrl"   json:"publisherUrl"`
	Glossaries   []Glossary    `bson:"glossaries"     json:"glossaries"`
}
