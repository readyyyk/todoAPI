package m_group

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"github.com/readyyyk/todoAPI/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"time"
)

func Create(c *gin.Context) {
	// parse user id
	uid, err := proceeding.ParseJWT(c.GetHeader("Auth"))
	if err != nil {
		c.JSON(401, apiErrors.Errors[4])
		logs.LogError(err)
		return
	}

	// create new m_group object
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		logs.LogError(err)
		return
	}
	var newGroup types.Group
	err = json.Unmarshal(jsonData, &newGroup)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		logs.LogError(err)
		return
	}

	newGroup.Id = primitive.NewObjectID()
	newGroup.Owners = []primitive.ObjectID{uid} // initial owner of m_group is user provided in jwt
	newGroup.Created = time.Now()

	// insert new m_group to database
	client := proceeding.NewDbClient()
	res, err := client.Database("todos").Collection("groups").InsertOne(ctx, newGroup)
	logs.LogError(err)
	c.JSON(http.StatusOK, res)
}
