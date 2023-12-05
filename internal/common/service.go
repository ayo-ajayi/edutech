package common

import (
	"github.com/ayo-ajayi/edutech/internal/student"
	"github.com/ayo-ajayi/edutech/internal/tutor"
	"go.mongodb.org/mongo-driver/bson"
)

type UserService struct {
	tutorRepo   tutor.ITutorRepo
	studentRepo student.IStudentRepo
}

func NewUserService(tutorRepo tutor.ITutorRepo, studentRepo student.IStudentRepo) *UserService {
	return &UserService{tutorRepo: tutorRepo, studentRepo: studentRepo}
}

func (us *UserService) GetTutors() ([]*tutor.Tutor, error) {
	return us.tutorRepo.GetTutors(bson.M{})
}

func (us *UserService) GetTutorsBySubjectId(subjectId string) ([]*tutor.Tutor, error) {
	return us.tutorRepo.GetTutors(bson.M{"subject": subjectId})
}

type IUserService interface {
	GetTutors() ([]*tutor.Tutor, error)
}
