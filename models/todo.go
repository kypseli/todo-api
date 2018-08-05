package models

import "gopkg.in/mgo.v2/bson"

// Represents a todo, we use bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type (
	Todo struct {
		ID          bson.ObjectId `bson:"_id" json:"id"`
		Description string        `bson:"description" json:"description"`
		Completed   bool          `bson:"completed" json:"completed"`
	}
)