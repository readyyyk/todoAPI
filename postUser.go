package main

import (
	"context"
	"encoding/base64"
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

func postUser(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	logs.LogError(err)
	//logs.AsJSON(jsonData)
	var newUser User
	err = json.Unmarshal(jsonData, &newUser)
	newUser.Id = primitive.NewObjectID()
	newUser.Registered = time.Now()
	newUser.Password = base64.StdEncoding.EncodeToString([]byte(newUser.Password))

	if validator.New().Struct(newUser) != nil || err != nil {
		logs.LogError(err)
		c.JSON(http.StatusBadRequest, map[string]string{"code": "0", "error": "Invalid data"})
		logs.LogError(errors.New(validator.New().Struct(newUser).Error()))
		return
	}

	if client.Database("todos").Collection("users").FindOne(context.TODO(), bson.D{{"email", newUser.Email}}).Err() != mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, map[string]string{"code": "1", "error": "User with this email already exists"})
		return
	}

	res, err := client.Database("todos").Collection("users").InsertOne(
		context.TODO(),
		newUser,
	)
	logs.LogError(err)

	c.JSON(http.StatusOK, res)
}
