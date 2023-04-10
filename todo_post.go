package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"os"
	"time"
)

func createTodo(c *gin.Context) {
	authToken := c.GetHeader("Auth")

	// parse and validate jwt auth token
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(authToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if !token.Valid || err != nil {
		c.JSON(401, errorDescriptionT{Code: 4, Description: "JWT token is invalid"})
		logs.LogError(err)
		return
	}

	// parse userId from jwt
	uid, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err == primitive.ErrInvalidHex {
		c.JSON(401, errorDescriptionT{Code: 4, Description: "JWT token is invalid"})
		logs.LogError(err)
		return
	}

	// check if current user owns this group
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	groupId, err := primitive.ObjectIDFromHex(c.Param("group_id"))
	if err == primitive.ErrInvalidHex {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "invalid data"})
		logs.LogError(err)
		return
	}
	groupFindRes := client.Database("todos").Collection("groups").FindOne(
		ctx,
		bson.D{{
			"$and",
			bson.A{
				bson.D{{"_id", groupId}},
				bson.D{{"owners", uid}},
			},
		}},
	)
	if groupFindRes.Err() == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, errorDescriptionT{
			Code:        5,
			Description: "Group don't exists",
		})
		return
	}

	// create new task object
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "Invalid data"})
		logs.LogError(err)
		return
	}
	var newTodo Todo
	err = json.Unmarshal(jsonData, &newTodo)

	newTodo.Id = primitive.NewObjectID()
	newTodo.Group = groupId
	newTodo.State = "passive"
	newTodo.StartDate = time.Now()

	if validator.New().Struct(newTodo) != nil || newTodo.Deadline.Before(time.Now()) || err != nil {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "Invalid data"})
		logs.LogError(errors.New(validator.New().Struct(newTodo).Error()))
		return
	}

	// insert new task to database
	res, err := client.Database("todos").Collection("todos").InsertOne(ctx, newTodo)
	logs.LogError(err)
	c.JSON(http.StatusOK, res)
}
