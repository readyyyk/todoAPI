package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"time"
)

func createTodo(c *gin.Context) {

	// parse user id
	uid, err := parseJWT(c.GetHeader("Auth"))
	if err != nil {
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
		return
	}
	logs.LogError(err)
	groupFindRes := client.Database("todos").Collection("groups").FindOne(
		ctx,
		/*bson.D{{
			"$and",
			bson.A{
				bson.D{{"_id", groupId}},
				bson.D{{"owners", uid}},
			},
		}},*/
		bson.D{{"_id", groupId}},
	)
	if groupFindRes.Err() == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, errorDescriptionT{Code: 5, Description: "Group don't exists"})
		return
	}
	var groupFound Group
	logs.LogError(groupFindRes.Decode(&groupFound))
	if !contains(groupFound.Owners, uid) {
		c.JSON(http.StatusForbidden, errorDescriptionT{Code: 6, Description: "User doesn't own this group"})
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
