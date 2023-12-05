package student

import (
	"errors"
	"time"

	"github.com/ayo-ajayi/edutech/internal/subject"
	"github.com/ayo-ajayi/edutech/internal/tutor"
	"github.com/ayo-ajayi/edutech/internal/user"
	"github.com/ayo-ajayi/edutech/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentService struct {
	verificationTokenManager utils.IVerificationTokenManager
	accessTokenManager       utils.IAccessTokenManager
	studentRepo              IStudentRepo
	emailManager             utils.IEmailManager
	subjectRepo              subject.IStudentSubjectRepo
	tutorRepo                tutor.IStudentTutorRepo
	studentSubjectTutorRepo  subject.IStudentSubjectTutorRepo
	baseUrl                  string
}

func NewStudentService(
	studentRepo IStudentRepo,
	verificationTokenManager utils.IVerificationTokenManager,
	accessTokenManager utils.IAccessTokenManager,
	emailManager utils.IEmailManager,
	subjectRepo subject.IStudentSubjectRepo,
	tutorRepo tutor.IStudentTutorRepo,
	studentSubjectTutorRepo subject.IStudentSubjectTutorRepo,
	baseUrl string,
) *StudentService {
	return &StudentService{studentRepo: studentRepo, verificationTokenManager: verificationTokenManager, accessTokenManager: accessTokenManager, emailManager: emailManager, baseUrl: baseUrl, subjectRepo: subjectRepo, tutorRepo: tutorRepo, studentSubjectTutorRepo: studentSubjectTutorRepo}
}

func (ss *StudentService) SignUpStudent(student *Student) error {
	exists, err := ss.studentRepo.StudentExists(bson.M{"email": student.Email})
	if err != nil {
		return err
	}
	if exists {
		return errors.New("student already exists")
	}
	passwordHash, err := utils.HashPassword(student.Password)
	if err != nil {
		return err
	}
	if passwordHash == "" {
		return errors.New("password hash is empty")
	}
	student.Password = passwordHash
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()
	student.Role = user.Student

	if err := ss.studentRepo.CreateStudent(student); err != nil {
		return err
	}

	verificationToken := utils.CreateVerificationToken()

	if err := ss.verificationTokenManager.SaveVerificationToken(student.Email, verificationToken); err != nil {
		return err
	}
	verificationLink, err := utils.ConstructVerificationLink(ss.baseUrl, "verify", verificationToken, student.Email)
	if err != nil {
		return err
	}
	if err := ss.emailManager.SendSignUpVerificationToken(student.Email, student.Firstname, verificationLink); err != nil {
		return err
	}
	return nil
}

func (ss *StudentService) GetStudent(id primitive.ObjectID) (*Student, error) {
	student, err := ss.studentRepo.GetStudent(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return student, nil
}
func (ss *StudentService) RegisterSubject(subjectId primitive.ObjectID, userId primitive.ObjectID) error {
	student, err := ss.studentRepo.GetStudent(bson.M{"_id": userId})
	if err != nil {
		return err
	}
	subject, err := ss.subjectRepo.GetSubject(bson.M{"_id": subjectId})
	if err != nil {
		return err
	}
	if subject == nil {
		return errors.New("subject does not exist")
	}
	for _, subject := range student.Subjects {
		if subject == subjectId {
			return errors.New("subject already registered")
		}
	}
	student.Subjects = append(student.Subjects, subjectId)
	return ss.studentRepo.UpdateStudent(bson.M{"_id": userId}, student)
}
func (ss *StudentService) GetRegisteredSubjects(userId primitive.ObjectID) ([]*subject.Subject, error) {
	student, err := ss.studentRepo.GetStudent(bson.M{"_id": userId})
	if err != nil {
		return nil, err
	}
	subjects := student.Subjects
	return ss.subjectRepo.GetSubjects(bson.M{"_id": bson.M{"$in": subjects}})
}

func (ss *StudentService) RegisterTutor(tutorId primitive.ObjectID, userId primitive.ObjectID) error {

	tutor, err := ss.tutorRepo.GetTutor(bson.M{"_id": tutorId})
	if err != nil {
		return err
	}
	subjects, err := ss.GetRegisteredSubjects(userId)
	if err != nil {
		return err
	}
	registered := false
	for _, s := range subjects {
		if tutor.Subject == s.Id {
			registered = true
			break
		}
	}
	if !registered {
		return errors.New("tutor's subject not registered by student")
	}
	exists, err := ss.studentSubjectTutorRepo.StudentSubjectTutorExists(bson.M{"student": userId, "tutor": tutorId})
	if err != nil {
		return err
	}
	if exists {
		return errors.New("tutor already registered")
	}
	studentSubjectTutor := &subject.StudentSubjectTutor{
		Id:        primitive.NewObjectID(),
		StudentId: userId,
		TutorId:   tutorId,
		SubjectId: tutor.Subject,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := ss.studentSubjectTutorRepo.CreateStudentSubjectTutor(studentSubjectTutor); err != nil {
		return err
	}
	return nil

}

func (ss *StudentService) GetRegisteredTutors(userId primitive.ObjectID) ([]*utils.StudentRegisteredTutorRes, error) {
	studentSubjectTutors, err := ss.studentSubjectTutorRepo.GetStudentSubjectTutors(bson.M{"student": userId})
	if err != nil {
		return nil, err
	}
	tutors := []*utils.StudentRegisteredTutorRes{}

	for _, studentSubjectTutor := range studentSubjectTutors {
		tutor, err := ss.tutorRepo.GetTutor(bson.M{"_id": studentSubjectTutor.TutorId})
		if err != nil {
			return nil, err
		}
		tutors = append(tutors, &utils.StudentRegisteredTutorRes{
			TutorId:      tutor.Id,
			SubjectId:    studentSubjectTutor.SubjectId,
			Email:        tutor.Email,
			FirstName:    tutor.Firstname,
			LastName:     tutor.Lastname,
			RegisteredAt: studentSubjectTutor.CreatedAt,
		})
	}
	return tutors, nil
}

type IStudentService interface {
	SignUpStudent(student *Student) error
	GetStudent(id primitive.ObjectID) (*Student, error)
	RegisterSubject(subjectId primitive.ObjectID, userId primitive.ObjectID) error
	GetRegisteredSubjects(userId primitive.ObjectID) ([]*subject.Subject, error)
	RegisterTutor(tutorId primitive.ObjectID, userId primitive.ObjectID) error
	GetRegisteredTutors(userId primitive.ObjectID) ([]*utils.StudentRegisteredTutorRes, error)
}
