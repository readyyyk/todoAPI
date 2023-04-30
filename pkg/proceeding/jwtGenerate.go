package proceeding

import (
	"github.com/golang-jwt/jwt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

func GenerateJWT(id primitive.ObjectID) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	logs.LogError(err)

	return signedToken
}
