package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
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

	oid, err := primitive.ObjectIDFromHex(c.Param("group_id"))
	if err == mongo.ErrInvalidIndexValue {
		c.JSON(http.StatusBadRequest, errorDescriptionT{
			Code:        0,
			Description: "Invalid data",
		})
		return
	}

	// check if current user owns this group
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var groupFind Group
	groupFindRes := client.Database("todos").Collection("groups").FindOne(
		ctx,
		bson.D{{
			"$and",
			bson.A{
				bson.D{{"_id", oid}},
				bson.D{{"owners", uid}},
			},
		}},
	)
	if groupFindRes.Err() == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, errorDescriptionT{
			Code:        5,
			Description: "Group don't exist",
		})
		return
	}
	logs.LogError(groupFindRes.Decode(&groupFind))
	if !contains(groupFind.Owners, uid) {
		c.JSON(http.StatusNotAcceptable, errorDescriptionT{
			Code:        6,
			Description: "User doesn't own this group",
		})
		return
	}

	delRes := make(map[string]int64)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	res, err := client.Database("todos").Collection("groups").DeleteOne(ctx, bson.D{{"_id", oid}})
	delRes["deletedGroupsCnt"] = res.DeletedCount
	logs.LogError(err)
	res, err = client.Database("todos").Collection("groups").DeleteMany(ctx, bson.D{{"group", oid}})
	delRes["deletedTodosCnt"] = res.DeletedCount
	logs.LogError(err)
	logs.AsJSON(delRes)
	c.JSON(http.StatusOK, delRes)
}
