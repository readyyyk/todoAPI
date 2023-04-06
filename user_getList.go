package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"time"
)

func getUserList(c *gin.Context) {
	if c.GetHeader("X-admin-access") != os.Getenv("ADMIN_ACCESS") {
		c.Status(http.StatusForbidden)
		return
	}

	users := client.Database("todos").Collection("users")
	var res []struct {
		Id    string `bson:"_id"`
		Name  string `bson:"name"`
		Email string `bson:"email"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logs.LogError(Select(users, ctx, bson.D{}, &res))

	c.JSON(http.StatusOK, res)
}
