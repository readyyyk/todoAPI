package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

func getUserData(c *gin.Context) /*(result []fullGroupData)*/ {
	type fullGroupData struct {
		GroupData Group
		TodosData []Todo
	}

	var result []fullGroupData

	cid := c.Param("id")
	id, err := primitive.ObjectIDFromHex(cid)
	if err == primitive.ErrInvalidHex {
		c.Status(http.StatusBadRequest)
		return
	}
	logs.LogError(err)

	// check if user owns provided group
	uid, err := parseJWT(c.GetHeader("Auth"))
	if err != nil {
		c.JSON(401, errorDescriptionT{Code: 4, Description: "JWT token is invalid"})
		logs.LogError(err)
		return
	}

	if uid != id && c.GetHeader("X-admin-access") != os.Getenv("ADMIN_ACCESS") {
		c.JSON(http.StatusForbidden, errorDescriptionT{Code: 6, Description: "User doesn't own this group"})
		return
	}

	// method logic
	groups := client.Database("todos").Collection("groups")
	todos := client.Database("todos").Collection("todos")

	var groupsData []Group
	logs.LogError(Select(groups, context.TODO(), bson.D{{"owners", id}}, &groupsData))

	for _, gr := range groupsData {
		var todosData []Todo
		logs.LogError(Select(todos, context.TODO(), bson.D{{"group", gr.Id}}, &todosData))

		result = append(result, fullGroupData{
			GroupData: gr,
			TodosData: todosData,
		})
	}

	c.IndentedJSON(http.StatusOK, result)
}
