package utils

import (
	"errors"

	"time"

	"github.com/ayo-ajayi/edutech/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccessDetails struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AccessUuid string             `json:"access_uuid" bson:"access_uuid"`
	UserId     primitive.ObjectID `json:"user_id" bson:"user_id"`
	ExpireAt   time.Time          `json:"expire_at" bson:"expire_at"`
}

type AccessTokenManager struct {
	accessTokenSecret          string
	accessTokenValidityInHours int64
	db                         db.IDatabase
}

type AccessTokenDetails struct {
	AccessToken string `json:"access_token"`
	AcessUuid   string `json:"-"`
	AtExpires   int64  `json:"at_expires"`
}

func NewTokenAccessManager(accessTokenSecret string, accessTokenValidityInHours int64, db db.IDatabase) *AccessTokenManager {
	return &AccessTokenManager{accessTokenSecret: accessTokenSecret, accessTokenValidityInHours: accessTokenValidityInHours, db: db}
}

func createAccessToken(userId primitive.ObjectID, uuid string, expires int64, secret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["access_uuid"] = uuid
	claims["exp"] = expires
	claims["authorized"] = true
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(secret))
}

func (atm *AccessTokenManager) GenerateAccessToken(userId primitive.ObjectID) (*AccessTokenDetails, error) {
	atd := &AccessTokenDetails{}
	atd.AtExpires = time.Now().Add(time.Hour * time.Duration(atm.accessTokenValidityInHours)).Unix()
	atd.AcessUuid = uuid.New().String()

	accessToken, err := createAccessToken(userId, atd.AcessUuid, atd.AtExpires, atm.accessTokenSecret)
	if err != nil {
		return nil, err
	}
	if accessToken == "" {
		return nil, errors.New("access token is empty")
	}
	atd.AccessToken = accessToken
	return atd, nil
}

func (atm *AccessTokenManager) SaveAccessToken(userId primitive.ObjectID, atd *AccessTokenDetails) error {
	exists, err := atm.accessTokenExists(userId)
	if err != nil {
		return err
	}
	if exists {
		if err := atm.DeleteAccessToken(bson.M{"user_id": userId}); err != nil {
			return err
		}
	}
	return atm.saveToken(userId, atd)
}
func (atm *AccessTokenManager) saveToken(userId primitive.ObjectID, atd *AccessTokenDetails) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := atm.db.InsertOne(ctx, &AccessDetails{
		AccessUuid: atd.AcessUuid,
		UserId:     userId,
		ExpireAt:   time.Unix(atd.AtExpires, 0),
	})
	return err
}

func (atm *AccessTokenManager) DeleteAccessToken(filter interface{}) error {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	_, err := atm.db.DeleteOne(ctx, filter)
	return err
}

func (atm *AccessTokenManager) accessTokenExists(userId primitive.ObjectID) (bool, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	err := atm.db.FindOne(ctx, bson.M{"user_id": userId}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (atm *AccessTokenManager) FindAccessToken(uuid string) (*AccessDetails, error) {
	ctx, cancel := db.DBReqContext(5)
	defer cancel()
	var accessDetails AccessDetails
	err := atm.db.FindOne(ctx, bson.M{"access_uuid": uuid}).Decode(&accessDetails)
	if err != nil {
		return nil, err
	}
	return &accessDetails, nil
}

func (atm *AccessTokenManager) ValidateAccessToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
}

func (atm *AccessTokenManager) ExtractAccessTokenMetadata(token *jwt.Token) (*AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("unauthorized")
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok || accessUuid == "" {
		return nil, errors.New("unauthorized")
	}

	userId, ok := claims["user_id"].(string)
	if !ok || userId == "" {
		return nil, errors.New("unauthorized")
	}
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	return &AccessDetails{
		AccessUuid: accessUuid,
		UserId:     userID,
	}, nil
}

type IMiddlewareAccessTokenManager interface {
	ExtractAccessTokenMetadata(token *jwt.Token) (*AccessDetails, error)
	ValidateAccessToken(token string, secret string) (*jwt.Token, error)
	FindAccessToken(uuid string) (*AccessDetails, error)
}

type IAccessTokenManager interface {
	SaveAccessToken(userId primitive.ObjectID, atd *AccessTokenDetails) error
	GenerateAccessToken(userId primitive.ObjectID) (*AccessTokenDetails, error)
	DeleteAccessToken(filter interface{}) error
}
