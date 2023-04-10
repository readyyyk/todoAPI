package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

	groups := client.Database("todos").Collection("groups")
	todos := client.Database("todos").Collection("todos")

	var groupsData []Group
	logs.LogError(Select(groups, context.TODO(), bson.D{{"owners", id}}, &groupsData))
	//fmt.Println(groupsData)

	for _, gr := range groupsData {
		var todosData []Todo
		logs.LogError(Select(todos, context.TODO(), bson.D{{"group", gr.Id}}, &todosData))

		result = append(result, fullGroupData{
			GroupData: gr,
			TodosData: todosData,
		})
		//logs.AsJSON(todosData)
	}

	c.IndentedJSON(http.StatusOK, result)
}
