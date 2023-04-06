package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"reflect"
)

func updateUser(c *gin.Context) {
	usersColl := client.Database("todos").Collection("users")

	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err == mongo.ErrInvalidIndexValue {
		c.JSON(http.StatusBadRequest,
			errorDescriptionT{
				Code:        0,
				Description: "Invalid data",
			})
		return
	} else if err != nil {
		logs.LogError(err)
	}

	reqBodyJSON, err := io.ReadAll(c.Request.Body)
	if len(reqBodyJSON) == 0 {
		c.Status(200)
		return
	}

	reqBody := struct {
		Email    string `bson:"email" json:"email" validate:"email"`
		Password string `bson:"password" json:"password" validate:"base64"` // base64
		Name     string `bson:"name" json:"name"`
	}{}
	logs.LogError(json.Unmarshal(reqBodyJSON, &reqBody))

	var currentUser User
	var selectRes []User
	logs.LogError(Select(usersColl, context.Background(), bson.D{{"_id", oid}}, &selectRes))
	if len(selectRes) == 0 {
		c.JSON(http.StatusNotFound, errorDescriptionT{
			Code:        2,
			Description: "user don't exists",
		})
		return
	}
	currentUser = selectRes[0]

	reflectCurrentUser := reflect.ValueOf(currentUser)
	reflectReqBody := reflect.ValueOf(reqBody)

	for i := 0; i < reflectReqBody.NumField(); i++ {
		if !reflectReqBody.Field(i).IsNil() {
			reflectCurrentUser.FieldByName(reflectReqBody.Type().Field(i).Name).Set(
				reflectReqBody.Field(i),
			)
		}
	}

	//logs.LogError(json.Unmarshal(reqBodyJSON, &currentUser))

	logs.AsJSON(reflectCurrentUser)

	//updRes, err := usersColl.ReplaceOne(context.Background(), bson.D{{"_id", oid}}, currentUser)
	logs.LogError(err)

	//fmt.Println(updRes)
}
