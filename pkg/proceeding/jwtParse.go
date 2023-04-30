package proceeding

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"regexp"
)

func ParseJWT(authToken string) (primitive.ObjectID, error) {
	// parse and validate jwt auth token
	if m, err := regexp.MatchString(`.+[.].+[.].+`, authToken); !m || err != nil {
		logs.LogError(err)
		return primitive.ObjectID{}, errors.New("JWT token is invalid")
	}

	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(authToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err == primitive.ErrInvalidHex {
		return primitive.ObjectID{}, err
	}
	if !token.Valid {
		return primitive.ObjectID{}, errors.New("JWT token is invalid")
	}

	// parse userId from jwt
	uid, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err == primitive.ErrInvalidHex {
		return primitive.ObjectID{}, err
	}

	return uid, nil
}
