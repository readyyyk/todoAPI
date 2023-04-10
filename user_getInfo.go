package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func getUserInfo(c *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err == primitive.ErrInvalidHex {
		c.Status(http.StatusBadRequest)
		return
	}
	logs.LogError(err)

	var currentUser []User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logs.LogError(Select(
		client.Database("todos").Collection("users"),
		ctx,
		bson.D{{
			"_id",
			userId,
		}},
		&currentUser,
	))

	if len(currentUser) == 0 {
		c.JSON(http.StatusNotFound, errorDescriptionT{
			Code:        2,
			Description: "user don't exists",
		})
		return
	}

	c.JSON(http.StatusOK, currentUser[0])
}
