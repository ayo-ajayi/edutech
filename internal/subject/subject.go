//subjects enrolled in by a student
//subjects taught by a tutor

package subject

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type Subject struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Compulsory bool               `json:"compulsory" bson:"compulsory"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

type StudentSubjectTutor struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentId primitive.ObjectID `json:"student_id" bson:"student_id"`
	SubjectId primitive.ObjectID `json:"subject_id" bson:"subject_id"`
	TutorId   primitive.ObjectID `json:"tutor_id" bson:"tutor_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

//multiple tutors for a subject
//timeline
