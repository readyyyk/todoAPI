package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func contains[T comparable](s []T, str T) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func deleteGroup(c *gin.Context) {
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

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, err = client.Database("todos").Collection("groups").DeleteOne(ctx, bson.D{{"_id", groupId}})
	logs.LogError(err)
	res, err := client.Database("todos").Collection("todos").DeleteMany(ctx, bson.D{{"group", groupId}})
	logs.LogError(err)
	c.JSON(http.StatusOK, res.DeletedCount)
}
