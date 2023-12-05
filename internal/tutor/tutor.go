package tutor

import (
	"github.com/ayo-ajayi/edutech/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tutor struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	*user.User
	Approved bool               `json:"approved" bson:"approved"`
	Subject  primitive.ObjectID `json:"subject" bson:"subject"`
}

//for now a tutor can only one course

//review, bio, rating, avatar
