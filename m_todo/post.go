package m_todo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"github.com/readyyyk/todoAPI/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"time"
)

/*
func contains[T comparable](s []T, str T) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
*/

func Create(c *gin.Context, client *mongo.Client) {
	// parse user id
	uid, err := proceeding.ParseJWT(c.GetHeader("Auth"))
	if err != nil {
		c.JSON(401, apiErrors.Errors[4])
		logs.LogError(err)
		return
	}

	// check if current user owns this m_group
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	groupId, err := primitive.ObjectIDFromHex(c.Param("group_id"))
	if err == primitive.ErrInvalidHex {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
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
		c.JSON(http.StatusNotFound, apiErrors.Errors[5])
		return
	}
	var groupFound types.Group
	logs.LogError(groupFindRes.Decode(&groupFound))
	if !proceeding.Contains(groupFound.Owners, uid) {
		c.JSON(http.StatusForbidden, apiErrors.Errors[6])
		return
	}

	// create new task object
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		logs.LogError(err)
		return
	}
	var newTodo types.Todo
	err = json.Unmarshal(jsonData, &newTodo)

	newTodo.Id = primitive.NewObjectID()
	newTodo.Group = groupId
	newTodo.State = "passive"
	newTodo.StartDate = time.Now()

	if validator.New().Struct(newTodo) != nil || newTodo.Deadline.Before(time.Now()) || err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		logs.LogError(errors.New(validator.New().Struct(newTodo).Error()))
		return
	}

	// insert new task to database
	res, err := client.Database("todos").Collection("todos").InsertOne(ctx, newTodo)
	logs.LogError(err)
	c.JSON(http.StatusOK, res)
}
