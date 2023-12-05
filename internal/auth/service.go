package auth

import (
	"errors"
	"time"

	"github.com/ayo-ajayi/edutech/internal/student"
	"github.com/ayo-ajayi/edutech/internal/subject"
	"github.com/ayo-ajayi/edutech/internal/tutor"
	"github.com/ayo-ajayi/edutech/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	baseUrl                  string
	emailManager             utils.IEmailManager
	accessTokenManager       utils.IAccessTokenManager
	verificationTokenManager utils.IVerificationTokenManager
	tutorRepo                tutor.ITutorRepo
	studentRepo              student.IStudentRepo
	subjectRepo              subject.IStudentSubjectRepo
}

func NewAuthService(tutorRepo tutor.ITutorRepo, studentRepo student.IStudentRepo, subjectRepo subject.ISubjectRepo, accessTokenManager utils.IAccessTokenManager, verificationTokenManager utils.IVerificationTokenManager, emailManager utils.IEmailManager, baseUrl string) *AuthService {
	return &AuthService{tutorRepo: tutorRepo,
		studentRepo:              studentRepo,
		accessTokenManager:       accessTokenManager,
		verificationTokenManager: verificationTokenManager,
		emailManager:             emailManager,
		baseUrl:                  baseUrl,
		subjectRepo:              subjectRepo,
	}
}

func (as *AuthService) Verify(email, token string) error {
	valid, err := as.verificationTokenManager.ValidateVerificationToken(email, token)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid token")
	}
	tutor, err := as.tutorRepo.GetTutor(bson.M{"user.email": email})
	if err == nil && tutor != nil {
		if err := as.tutorRepo.UpdateTutor(bson.M{"user.email": email}, bson.M{"$set": bson.M{"user.is_verified": true}}); err != nil {
			return err
		}
		return nil
	}

	student, err := as.studentRepo.GetStudent(bson.M{"user.email": email})
	if err == nil && student != nil {
		subjects, err := as.subjectRepo.GetSubjects(bson.M{"compulsory": true})
		if err != nil {
			return err
		}
		compulsorySubjects := []primitive.ObjectID{}
		for _, subject := range subjects {
			compulsorySubjects = append(compulsorySubjects, subject.Id)
		}
		if err := as.studentRepo.UpdateStudent(bson.M{"user.email": email}, bson.M{"$set": bson.M{"user.is_verified": true, "subjects": compulsorySubjects}}); err != nil {
			return err
		}
		//delete token after user is updated
		return nil
	}
	return errors.New("invalid email")
}

func (as *AuthService) Login(email, password string) (interface{}, *utils.AccessTokenDetails, error) {
	tutor, err := as.tutorRepo.GetTutor(bson.M{"user.email": email})
	if err == nil && tutor != nil {
		if !tutor.IsVerified {
			return nil, nil, errors.New("tutor not verified")
		}
		if !utils.CheckPasswordHash(password, tutor.Password) {
			return nil, nil, errors.New("invalid username or password")
		}
		accessTokenDetails, err := as.accessToken(tutor.Id)
		if err != nil {
			return nil, nil, err
		}
		return tutor, accessTokenDetails, nil
	}

	student, err := as.studentRepo.GetStudent(bson.M{"user.email": email})
	if err == nil && student != nil {
		if !student.IsVerified {
			return nil, nil, errors.New("student not verified")
		}
		if !utils.CheckPasswordHash(password, student.Password) {
			return nil, nil, errors.New("invalid username or password")
		}
		accessTokenDetails, err := as.accessToken(student.Id)
		if err != nil {
			return nil, nil, err
		}
		return student, accessTokenDetails, nil
	}
	return nil, nil, errors.New("invalid username or password")
}

func (as *AuthService) accessToken(id primitive.ObjectID) (*utils.AccessTokenDetails, error) {
	accessToken, err := as.accessTokenManager.GenerateAccessToken(id)
	if err != nil {
		return nil, err
	}
	if err := as.accessTokenManager.SaveAccessToken(id, accessToken); err != nil {
		return nil, err
	}
	if accessToken == nil {
		return nil, errors.New("access token is empty")
	}
	return &utils.AccessTokenDetails{
		AccessToken: accessToken.AccessToken,
		AtExpires:   accessToken.AtExpires,
	}, nil
}

func (as *AuthService) Logout(accessUuid string) error {
	return as.accessTokenManager.DeleteAccessToken(bson.M{"access_uuid": accessUuid})
}
func (as *AuthService) ForgotPassword(email string) error {
	resetToken := utils.CreateVerificationToken()
	if err := as.verificationTokenManager.SaveVerificationToken(email, resetToken); err != nil {
		return err
	}
	link, err := utils.ConstructVerificationLink(as.baseUrl, "verify-reset-token", resetToken, email)
	if err != nil {
		return err
	}
	tutor, err := as.tutorRepo.GetTutor(bson.M{"user.email": email})
	if err == nil && tutor != nil {
		if err := as.emailManager.SendResetPasswordToken(email, tutor.Firstname, link); err != nil {
			return err
		}
		return nil
	}
	student, err := as.studentRepo.GetStudent(bson.M{"user.email": email})
	if err == nil && student != nil {
		if err := as.emailManager.SendResetPasswordToken(email, student.Firstname, link); err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid email")

}

func (as *AuthService) ResetPassword(email, password, token string) error {
	valid, err := as.verificationTokenManager.ValidateVerificationToken(email, token)
	if err != nil || !valid {
		return errors.New("invalid or expired token")
	}
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	if passwordHash == "" {
		return errors.New("password hash is empty")
	}
	tutor, err := as.tutorRepo.GetTutor(bson.M{"user.email": email})
	if err == nil && tutor != nil {
		if err := as.tutorRepo.UpdateTutor(bson.M{"user.email": email}, bson.M{"$set": bson.M{"user.password": passwordHash, "user.updated_at": time.Now()}}); err != nil {
			return err
		}
		return nil
	}

	student, err := as.studentRepo.GetStudent(bson.M{"user.email": email})
	if err == nil && student != nil {
		if err := as.studentRepo.UpdateStudent(bson.M{"user.email": email}, bson.M{"$set": bson.M{"user.password": passwordHash, "user.updated_at": time.Now()}}); err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid email")
}

type IAuthService interface {
	Verify(email, token string) error
	Login(email, password string) (interface{}, *utils.AccessTokenDetails, error)
	Logout(accessUuid string) error
	ForgotPassword(email string) error
	ResetPassword(email, password, token string) error
}
