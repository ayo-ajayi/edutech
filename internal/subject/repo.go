package subject

import (
	"github.com/ayo-ajayi/edutech/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubjectRepo struct {
	db db.IDatabase
}

func NewSubjectRepo(db db.IDatabase) *SubjectRepo {
	return &SubjectRepo{db: db}
}

func (sr *SubjectRepo) CreateSubject(subject *Subject) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := sr.db.InsertOne(ctx, subject)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubjectRepo) CreateSubjects(subjects []*Subject) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var docs []interface{}
	for _, subject := range subjects {
		docs = append(docs, subject)
	}
	_, err := sr.db.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubjectRepo) SubjectExists(filter interface{}) (bool, error) {
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

func (sr *SubjectRepo) UpdateSubject(filter interface{}, update interface{}) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := sr.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubjectRepo) GetSubject(filter interface{}) (*Subject, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var subject Subject
	err := sr.db.FindOne(ctx, filter).Decode(&subject)
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func (sr *SubjectRepo) GetSubjects(filter interface{}) ([]*Subject, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var subjects []*Subject
	cursor, err := sr.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &subjects)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}

type ISubjectRepo interface {
	CreateSubject(subject *Subject) error
	CreateSubjects(subjects []*Subject) error
	SubjectExists(filter interface{}) (bool, error)
	UpdateSubject(filter interface{}, update interface{}) error
	GetSubject(filter interface{}) (*Subject, error)
	GetSubjects(filter interface{}) ([]*Subject, error)
}

type IStudentSubjectRepo interface {
	GetSubjects(filter interface{}) ([]*Subject, error)
	GetSubject(filter interface{}) (*Subject, error)
}

type StudentSubjectTutorRepo struct {
	db db.IDatabase
}

func NewStudentSubjectTutorRepo(db db.IDatabase) *StudentSubjectTutorRepo {
	return &StudentSubjectTutorRepo{db: db}
}

func (sstr *StudentSubjectTutorRepo) CreateStudentSubjectTutor(studentSubjectTutor *StudentSubjectTutor) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := sstr.db.InsertOne(ctx, studentSubjectTutor)
	if err != nil {
		return err
	}
	return nil
}

func (sstr *StudentSubjectTutorRepo) GetStudentSubjectTutor(filter interface{}) (*StudentSubjectTutor, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var studentSubjectTutor StudentSubjectTutor
	err := sstr.db.FindOne(ctx, filter).Decode(&studentSubjectTutor)
	if err != nil {
		return nil, err
	}
	return &studentSubjectTutor, nil
}

func (sstr *StudentSubjectTutorRepo) GetStudentSubjectTutors(filter interface{}) ([]*StudentSubjectTutor, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var studentSubjectTutors []*StudentSubjectTutor
	cursor, err := sstr.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &studentSubjectTutors)
	if err != nil {
		return nil, err
	}
	return studentSubjectTutors, nil
}

func (sstr *StudentSubjectTutorRepo) StudentSubjectTutorExists(filter interface{}) (bool, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	err := sstr.db.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type IStudentSubjectTutorRepo interface {
	CreateStudentSubjectTutor(studentSubjectTutor *StudentSubjectTutor) error
	StudentSubjectTutorExists(filter interface{}) (bool, error)
	GetStudentSubjectTutors(filter interface{}) ([]*StudentSubjectTutor, error)
}
