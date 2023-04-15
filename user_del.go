package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

// not completely tested
func deleteUser(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err == mongo.ErrInvalidIndexValue {
		c.JSON(http.StatusBadRequest, errorDescriptionT{Code: 0, Description: "Invalid data"})
		return
	}

	// check if user owns provided group
	uid, err := parseJWT(c.GetHeader("Auth"))
	if err != nil {
		c.JSON(401, errorDescriptionT{Code: 4, Description: "JWT token is invalid"})
		logs.LogError(err)
		return
	}

	if uid != oid && c.GetHeader("X-admin-access") != os.Getenv("ADMIN_ACCESS") {
		c.JSON(http.StatusForbidden, errorDescriptionT{Code: 6, Description: "User doesn't own this group"})
		return
	}

	// method logic
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// deleting user
	_, err = client.Database("todos").Collection("users").DeleteOne(ctx, bson.D{{"_id", oid}})
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, errorDescriptionT{
			Code:        2,
			Description: "user dont exists",
		})
		return
	}
	logs.LogError(err)

	// deleting all groups that are owned ONLY by this user
	delRes := make(map[string]int64)
	delRes["deletedTodosCnt"] = 0

	var groupIds []struct {
		Id primitive.ObjectID `bson:"_id" json:"id"`
	}
	logs.LogError(Select(client.Database("todos").Collection("groups"), ctx, bson.D{{"owners", bson.A{oid}}}, &groupIds))
	currentRes, err := client.Database("todos").Collection("groups").DeleteMany(ctx, bson.D{{"owners", bson.A{oid}}})
	logs.LogError(err)
	delRes["deletedGroupsCnt"] = currentRes.DeletedCount

	// deleting todos
	for _, groupId := range groupIds {
		currentRes, err := client.Database("todos").Collection("groups").DeleteMany(ctx, bson.D{{"group", groupId}})
		logs.LogError(err)
		delRes["deletedTodosCnt"] += currentRes.DeletedCount
	}

	c.JSON(http.StatusOK, delRes)
}
