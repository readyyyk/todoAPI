package m_user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"github.com/readyyyk/todoAPI/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetInfo(c *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err == primitive.ErrInvalidHex {
		c.Status(http.StatusBadRequest)
		return
	}
	logs.LogError(err)

	var currentUser []types.User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := proceeding.NewDbClient()
	logs.LogError(proceeding.Select(
		client.Database("todos").Collection("users"),
		ctx,
		bson.D{{
			"_id",
			userId,
		}},
		&currentUser,
	))

	if len(currentUser) == 0 {
		c.JSON(http.StatusNotFound, apiErrors.Errors[2])
		return
	}

	c.JSON(http.StatusOK, currentUser[0])
}
