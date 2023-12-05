package subject

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubjectService struct {
	subjectRepo ISubjectRepo
}

func NewSubjectService(subjectRepo ISubjectRepo, compulsorySubjects ...string) (*SubjectService, error) {
	subjects := []*Subject{}
	for _, subject := range compulsorySubjects {
		exists, err := subjectRepo.SubjectExists(bson.M{"name": subject})
		if err != nil {
			return nil, err
		}
		if !exists {
			subjects = append(subjects, &Subject{
				Id:         primitive.NewObjectID(),
				Name:       subject,
				Compulsory: true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}
	}
	if len(subjects) > 0 {
		err := subjectRepo.CreateSubjects(subjects)
		if err != nil {
			return nil, err
		}
	}
	return &SubjectService{subjectRepo: subjectRepo}, nil
}

func (ss *SubjectService) CreateSubject(subject *Subject) error {
	subject.Id = primitive.NewObjectID()
	subject.CreatedAt = time.Now()
	subject.UpdatedAt = time.Now()

	err := ss.subjectRepo.CreateSubject(subject)
	if err != nil {
		return err
	}
	return nil
}

type ISubjectService interface {
	CreateSubject(subject *Subject) error
}
