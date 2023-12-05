package tutor

import (
	"github.com/ayo-ajayi/edutech/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type TutorRepo struct {
	db db.IDatabase
}

func NewTutorRepo(db db.IDatabase) *TutorRepo {
	return &TutorRepo{db: db}
}

func (tr *TutorRepo) CreateTutor(tutor *Tutor) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := tr.db.InsertOne(ctx, tutor)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TutorRepo) TutorExists(filter interface{}) (bool, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	err := tr.db.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (tr *TutorRepo) UpdateTutor(filter interface{}, update interface{}) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := tr.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TutorRepo) GetTutor(filter interface{}) (*Tutor, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var tutor Tutor
	err := tr.db.FindOne(ctx, filter).Decode(&tutor)
	if err != nil {
		return nil, err
	}
	return &tutor, nil
}

func (sr *TutorRepo) GetTutors(filter interface{}) ([]*Tutor, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var tutors []*Tutor
	cursor, err := sr.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &tutors)
	if err != nil {
		return nil, err
	}
	return tutors, nil
}

type ITutorRepo interface {
	CreateTutor(user *Tutor) error
	TutorExists(filter interface{}) (bool, error)
	UpdateTutor(filter interface{}, update interface{}) error
	GetTutor(filter interface{}) (*Tutor, error)
	GetTutors(filter interface{}) ([]*Tutor, error)
}

type IStudentTutorRepo interface {
	GetTutors(filter interface{}) ([]*Tutor, error)
	GetTutor(filter interface{}) (*Tutor, error)
}

type IMiddlewareTutorRepo interface {
	GetTutor(filter interface{}) (*Tutor, error)
}
