package student

import (
	"github.com/ayo-ajayi/edutech/internal/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	*user.User
	Subjects []primitive.ObjectID `json:"subjects" bson:"subjects"`
}
