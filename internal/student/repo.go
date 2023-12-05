package student

import (
	"github.com/ayo-ajayi/edutech/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentRepo struct {
	db db.IDatabase
}

func NewStudentRepo(db db.IDatabase) *StudentRepo {
	return &StudentRepo{db: db}
}

func (sr *StudentRepo) CreateStudent(Student *Student) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := sr.db.InsertOne(ctx, Student)
	if err != nil {
		return err
	}
	return nil
}

func (sr *StudentRepo) StudentExists(filter interface{}) (bool, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	err := sr.db.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (sr *StudentRepo) UpdateStudent(filter interface{}, update interface{}) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := sr.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (sr *StudentRepo) GetStudent(filter interface{}) (*Student, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var student Student
	err := sr.db.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

type IStudentRepo interface {
	CreateStudent(student *Student) error
	StudentExists(filter interface{}) (bool, error)
	UpdateStudent(filter interface{}, update interface{}) error
	GetStudent(filter interface{}) (*Student, error)
}

type IMiddlewareStudentRepo interface {
	GetStudent(filter interface{}) (*Student, error)
}
