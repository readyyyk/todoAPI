package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"time"
)

func createGroup(c *gin.Context) {

	// parse user id
	uid, err := parseJWT(c.GetHeader("Auth"))
	if err != nil {
		c.JSON(401, errorDescriptionT{Code: 4, Description: "JWT token is invalid"})
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
