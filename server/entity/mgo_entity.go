package entity

import "gopkg.in/mgo.v2/bson"

type PersonMG struct {
	Id		bson.ObjectId `bson:"_id"`
	Name 	string `bson:"name"`
	Age		int32 `bson:"age"`
}
