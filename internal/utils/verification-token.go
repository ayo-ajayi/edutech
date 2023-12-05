package utils

import (
	"errors"
	"time"

	"github.com/ayo-ajayi/edutech/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VerificationToken struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Token     string             `json:"token" bson:"token"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type VerificationTokenManager struct {
	db                                db.IDatabase
	SignUpTokenValidityInSecs         uint
	ForgotPasswordTokenValidityInSecs uint
}

type IVerificationTokenManager interface {
	SaveVerificationToken(email string, verificationToken string) error
	ValidateVerificationToken(email string, verificationToken string) (bool, error)
}

func NewVerificationTokenManager(db db.IDatabase, signUpTokenValidityInSecs uint, forgotPasswordTokenValidityInSecs uint) *VerificationTokenManager {
	return &VerificationTokenManager{db: db, SignUpTokenValidityInSecs: signUpTokenValidityInSecs, ForgotPasswordTokenValidityInSecs: forgotPasswordTokenValidityInSecs}
}

func InitVerificationTokenExpiryIndex(collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.M{"expires_at": 1}, Options: options.Index().SetExpireAfterSeconds(0),
	}
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return errors.New("Error creating TTL index for token collection:" + err.Error())
	}
	return nil
}

func (vtm *VerificationTokenManager) SaveVerificationToken(email string, verificationToken string) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	hashToken, err := HashPassword(verificationToken)
	if err != nil {
		return err
	}
	token := VerificationToken{
		Email:     email,
		Token:     hashToken,
		ExpiresAt: time.Now().Add(time.Duration(vtm.SignUpTokenValidityInSecs) * time.Second),
		CreatedAt: time.Now(),
	}
	if _, err := vtm.db.InsertOne(ctx, token); err != nil {
		return err
	}
	return nil
}
func (vtm *VerificationTokenManager) ValidateVerificationToken(email string, verificationToken string) (bool, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var token VerificationToken
	findOptions := options.FindOne().SetSort(bson.D{
		primitive.E{Key: "created_at", Value: -1},
	})
	err := vtm.db.FindOne(ctx, bson.M{"email": email}, findOptions).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	valid := CheckPasswordHash(verificationToken, token.Token)
	if !valid {
		return false, nil
	}
	return true, nil
}
