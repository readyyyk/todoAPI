package m_user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"time"
)

func Create(c *gin.Context, client *mongo.Client) {
	jsonData, err := io.ReadAll(c.Request.Body)
	logs.LogError(err)
	var newUser types.User
	err = json.Unmarshal(jsonData, &newUser)

	newUser.Id = primitive.NewObjectID()
	newUser.Registered = time.Now()
	newUser.Password = base64.StdEncoding.EncodeToString([]byte(newUser.Password))

	if validator.New().Struct(newUser) != nil || err != nil {
		logs.LogError(err)
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		//logs.LogError(errors.New(validator.New().Struct(newUser).Error()))
		return
	}

	if client.Database("todos").Collection("users").FindOne(context.TODO(), bson.D{{"email", newUser.Email}}).Err() != mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[1])
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.Database("todos").Collection("users").InsertOne(
		ctx,
		newUser,
	)
	logs.LogError(err)

	c.JSON(http.StatusOK, res)
}
