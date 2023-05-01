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
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
)

func GetData(c *gin.Context, client *mongo.Client) /*(result []fullGroupData)*/ {
	type fullGroupData struct {
		GroupData types.Group
		TodosData []types.Todo
	}

	var result []fullGroupData

	cid := c.Param("id")
	id, err := primitive.ObjectIDFromHex(cid)
	if err == primitive.ErrInvalidHex {
		c.Status(http.StatusBadRequest)
		return
	}
	logs.LogError(err)

	// check if user owns provided m_group
	uid, err := proceeding.ParseJWT(c.GetHeader("Auth"))
	if err != nil {
		c.JSON(401, apiErrors.Errors[4])
		logs.LogError(err)
		return
	}

	if uid != id && c.GetHeader("X-admin-access") != os.Getenv("ADMIN_ACCESS") {
		c.JSON(http.StatusForbidden, apiErrors.Errors[6])
		return
	}

	// method logic
	groups := client.Database("todos").Collection("groups")
	todos := client.Database("todos").Collection("todos")

	var groupsData []types.Group
	logs.LogError(proceeding.Select(groups, context.TODO(), bson.D{{"owners", id}}, &groupsData))

	for _, gr := range groupsData {
		var todosData []types.Todo
		logs.LogError(proceeding.Select(todos, context.TODO(), bson.D{{"m_group", gr.Id}}, &todosData))

		result = append(result, fullGroupData{
			GroupData: gr,
			TodosData: todosData,
		})
	}

	c.IndentedJSON(http.StatusOK, result)
}
