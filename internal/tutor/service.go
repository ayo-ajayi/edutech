package tutor

import (
	"errors"
	"time"

	"github.com/ayo-ajayi/edutech/internal/user"
	"github.com/ayo-ajayi/edutech/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TutorService struct {
	verificationTokenManager utils.IVerificationTokenManager
	accessTokenManager       utils.IAccessTokenManager
	tutorRepo                ITutorRepo
	emailManager             utils.IEmailManager
	baseUrl                  string
}

func NewTutorService(tutorRepo ITutorRepo, verificationTokenManager utils.IVerificationTokenManager, accessTokenManager utils.IAccessTokenManager, emailManager utils.IEmailManager, baseUrl string) *TutorService {
	return &TutorService{tutorRepo: tutorRepo, verificationTokenManager: verificationTokenManager, accessTokenManager: accessTokenManager, emailManager: emailManager, baseUrl: baseUrl}
}

func (ts *TutorService) SignUpTutor(tutor *Tutor) error {
	exists, err := ts.tutorRepo.TutorExists(bson.M{"email": tutor.Email})
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}
	passwordHash, err := utils.HashPassword(tutor.Password)
	if err != nil {
		return err
	}
	if passwordHash == "" {
		return errors.New("password hash is empty")
	}
	tutor.Password = passwordHash
	tutor.CreatedAt = time.Now()
	tutor.UpdatedAt = time.Now()
	tutor.Role = user.Tutor

	if err := ts.tutorRepo.CreateTutor(tutor); err != nil {
		return err
	}

	verificationToken := utils.CreateVerificationToken()

	if err := ts.verificationTokenManager.SaveVerificationToken(tutor.Email, verificationToken); err != nil {
		return err
	}
	verificationLink, err := utils.ConstructVerificationLink(ts.baseUrl, "verify", verificationToken, tutor.Email)
	if err != nil {
		return err
	}
	if err := ts.emailManager.SendSignUpVerificationToken(tutor.Email, tutor.Firstname, verificationLink); err != nil {
		return err
	}
	return nil
}

func (ts *TutorService) GetTutor(id primitive.ObjectID) (*Tutor, error) {
	tutor, err := ts.tutorRepo.GetTutor(bson.M{"_id": id})

	if err != nil {
		return nil, err
	}
	return tutor, nil
}

type ITutorService interface {
	SignUpTutor(tutor *Tutor) error
	GetTutor(id primitive.ObjectID) (*Tutor, error)
}
