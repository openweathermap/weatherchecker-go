package models

import "gopkg.in/mgo.v2/bson"

type DbEntryBase struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"objectid"`
}
