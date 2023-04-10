package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"time"
)

func createGroup(c *gin.Context) {
	authToken := c.GetHeader("Auth")

	// parse and validate jwt auth token
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(authToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if !token.Valid {
		c.JSON(401, errorDescriptionT{
			Code:        4,
			Description: "JWT token is invalid",
		})
		return
	}

	// parse userId from jwt
	uid, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err == primitive.ErrInvalidHex {
		c.JSON(401, errorDescriptionT{
			Code:        4,
			Description: "JWT token is invalid",
		})
		logs.LogError(err)
		return
	}

	// create new group object
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "Invalid data"})
		logs.LogError(err)
		return
	}
	var newGroup Group
	err = json.Unmarshal(jsonData, &newGroup)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "Invalid data"})
		logs.LogError(err)
		return
	}

	newGroup.Id = primitive.NewObjectID()
	newGroup.Owners = []primitive.ObjectID{uid} // initial owner of group is user provided in jwt
	newGroup.Created = time.Now()

	// insert new group to database
	res, err := client.Database("todos").Collection("groups").InsertOne(ctx, newGroup)
	logs.LogError(err)
	c.JSON(http.StatusOK, res)
}
