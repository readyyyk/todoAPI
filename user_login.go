package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
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

func loginUser(c *gin.Context) {

	type loginRespT struct {
		Logged bool              `json:"logged"`
		Err    errorDescriptionT `json:"err"`
		Token  string            `json:"token"`
	}
	type userLoginT struct {
		Email string `json:"email" validate:"required,email"`
		Pswd  string `json:"password" validate:"required,base64"`
	}

	var userData userLoginT
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "Invalid data"})
		return
	}

	err = json.Unmarshal(jsonData, &userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "Invalid data"})
		return
	}

	uemail := userData.Email
	upswd, err := base64.StdEncoding.DecodeString(userData.Pswd)

	// validate entered data
	if _, ok := err.(base64.CorruptInputError); ok || len(uemail) == 0 || len(upswd) == 0 {
		c.JSON(http.StatusBadRequest, loginRespT{
			false,
			errorDescriptionT{
				Code:        0,
				Description: "Invalid data",
			},
			"",
		})
		return
	} /*else if len(uemail) == 0 || len(upswd) == 0 {
		c.JSON(http.StatusBadRequest, errorDescriptionT{
			Code:        0,
			Description: "Cannot get username or password",
		})
		return
	}*/

	// select user with provided email
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	var userSelectRes []User
	logs.LogError(Select(
		client.Database("todos").Collection("users"),
		ctx,
		bson.D{{"email", uemail}},
		&userSelectRes,
	))
	if len(userSelectRes) == 0 {
		c.JSON(http.StatusNotFound, loginRespT{
			false,
			errorDescriptionT{
				Code:        2,
				Description: "user don't exists",
			},
			"",
		})
		return
	}
	userSelected := userSelectRes[0]

	// validate password
	if userSelected.Password != base64.StdEncoding.EncodeToString(upswd) {
		c.JSON(http.StatusNotFound, loginRespT{
			false,
			errorDescriptionT{
				Code:        3,
				Description: "Wrong password",
			},
			"",
		})
		return
	}

	// generate jwt
	signedToken := GenerateJWT(userSelected.Id)

	c.JSON(http.StatusOK, loginRespT{
		true,
		errorDescriptionT{},
		signedToken,
	})
}
