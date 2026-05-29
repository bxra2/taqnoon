package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Term struct {
	ID          bson.ObjectID `bson:"_id,omitempty"    json:"id"`
	English     string        `bson:"english"          json:"english"`
	Arabic      string        `bson:"arabic"           json:"arabic"`
	French      string        `bson:"french,omitempty" json:"french,omitempty"`
	German      string        `bson:"german,omitempty" json:"german,omitempty"`
	TURL        string        `bson:"tURL,omitempty"   json:"tURL,omitempty"`
	GlossaryEn  string        `bson:"glossaryEn"       json:"glossaryEn"`
	GlossaryAr  string        `bson:"glossaryAr"       json:"glossaryAr"`
	GlossaryURL string        `bson:"glossaryUrl"      json:"glossaryUrl"`
	PublisherID bson.ObjectID `bson:"publisherId"      json:"publisherId"`
}
