package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
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
	logs.LogError(Select(users, context.TODO(), bson.D{}, &res))

	c.JSON(http.StatusOK, res)
}
