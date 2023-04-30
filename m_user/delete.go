package m_user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

func Delete(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err == mongo.ErrInvalidIndexValue {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		return
	}

	// check if user owns provided m_group
	uid, err := proceeding.ParseJWT(c.GetHeader("Auth"))
	if err != nil {
		c.JSON(401, apiErrors.Errors[4])
		logs.LogError(err)
		return
	}

	if uid != oid && c.GetHeader("X-admin-access") != os.Getenv("ADMIN_ACCESS") {
		c.JSON(http.StatusForbidden, apiErrors.Errors[6])
		return
	}

	// method logic
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// deleting user
	client := proceeding.NewDbClient()
	_, err = client.Database("todos").Collection("users").DeleteOne(ctx, bson.D{{"_id", oid}})
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, apiErrors.Errors[2])
		return
	}
	logs.LogError(err)

	// deleting all groups that are owned ONLY by this user
	delRes := make(map[string]int64)
	delRes["deletedTodosCnt"] = 0

	var groupIds []struct {
		Id primitive.ObjectID `bson:"_id" json:"id"`
	}
	logs.LogError(proceeding.Select(client.Database("todos").Collection("groups"), ctx, bson.D{{"owners", bson.A{oid}}}, &groupIds))
	currentRes, err := client.Database("todos").Collection("groups").DeleteMany(ctx, bson.D{{"owners", bson.A{oid}}})
	logs.LogError(err)
	delRes["deletedGroupsCnt"] = currentRes.DeletedCount

	// deleting todos
	for _, groupId := range groupIds {
		currentRes, err := client.Database("todos").Collection("groups").DeleteMany(ctx, bson.D{{"m_group", groupId}})
		logs.LogError(err)
		delRes["deletedTodosCnt"] += currentRes.DeletedCount
	}

	c.JSON(http.StatusOK, delRes)
}
