package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ApptDate     time.Time          `bson:"apptDate" validate:"required"`
	User         primitive.ObjectID `bson:"user" validate:"required"`
	Hotel        primitive.ObjectID `bson:"hotel" validate:"required"`
	WifiPassword string             `bson:"wifiPassword,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty"`
}
