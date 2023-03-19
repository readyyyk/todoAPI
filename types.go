package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id         primitive.ObjectID `bson:"_id" json:"oid"`
	Email      string             `bson:"email" json:"email" validate:"required,email"`
	Password   string             `bson:"password" json:"password" validate:"required,base64"` // base64
	Registered time.Time          `bson:"registered" json:"registered"`
	Name       string             `bson:"name" json:"name" validate:"required"`
}

type Group struct {
	Id      primitive.ObjectID   `bson:"_id" json:"id"`
	Owners  []primitive.ObjectID `bson:"owners" json:"owners" validate:"required"` // foreign
	Title   string               `bson:"title" json:"title" validate:"required"`
	Created time.Time            `bson:"created" json:"created"`
}

type Todo struct {
	Id        primitive.ObjectID `bson:"_id" json:"oid"`
	Group     primitive.ObjectID `bson:"group" json:"group" validate:"required"` // foreign
	Title     string             `bson:"title" json:"title" validate:"required"`
	Text      string             `bson:"text" json:"text" validate:"required"`
	State     string             `bson:"state" json:"state" validate:"required"`
	StartDate time.Time          `bson:"startDate" json:"startDate"`
	Deadline  time.Time          `bson:"deadline" json:"deadline" validate:"required"`
}
