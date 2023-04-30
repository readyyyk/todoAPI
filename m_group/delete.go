package m_group

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"github.com/readyyyk/todoAPI/pkg/types"
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

func Delete(c *gin.Context) {
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
	client := proceeding.NewDbClient()
	groupFindRes := client.Database("todos").Collection("groups").FindOne(
		ctx,
		bson.D{{"_id", groupId}},
	)
	if groupFindRes.Err() == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, apiErrors.Errors[5])
		return
	}
	var groupFound types.Group
	logs.LogError(groupFindRes.Decode(&groupFound))
	if !contains(groupFound.Owners, uid) {
		c.JSON(http.StatusForbidden, apiErrors.Errors[6])
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, err = client.Database("todos").Collection("groups").DeleteOne(ctx, bson.D{{"_id", groupId}})
	logs.LogError(err)
	res, err := client.Database("todos").Collection("todos").DeleteMany(ctx, bson.D{{"m_group", groupId}})
	logs.LogError(err)
	c.JSON(http.StatusOK, res)
}
