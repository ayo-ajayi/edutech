package utils

import (
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConstructVerificationLink(baseUrl, path, verificationToken string, email string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	u.Path += "/" + path + "/" + verificationToken
	q := u.Query()
	q.Set("email", email)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func CreateVerificationToken() string {
	return primitive.NewObjectID().Hex()
}
