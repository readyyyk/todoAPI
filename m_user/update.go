package m_user

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"github.com/readyyyk/todoAPI/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"os"
	"reflect"
)

func Update(c *gin.Context) {
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err == mongo.ErrInvalidIndexValue {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		return
	} else if err != nil {
		logs.LogError(err)
	}

	// check if user tries to update itself or not
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
	reqBodyJSON, err := io.ReadAll(c.Request.Body)
	if len(reqBodyJSON) == 0 {
		c.Status(200)
		return
	}

	reqBody := struct {
		Email    string `bson:"email" json:"email" validate:"email"`
		Password string `bson:"password" json:"password" validate:"base64"` // base64
		Name     string `bson:"name" json:"name"`
	}{}
	logs.LogError(json.Unmarshal(reqBodyJSON, &reqBody))

	client := proceeding.NewDbClient()
	usersColl := client.Database("todos").Collection("users")
	
	var currentUser types.User
	var selectRes []types.User
	logs.LogError(proceeding.Select(usersColl, context.Background(), bson.D{{"_id", oid}}, &selectRes))
	if len(selectRes) == 0 {
		c.JSON(http.StatusNotFound, apiErrors.Errors[2])
		return
	}
	currentUser = selectRes[0]

	reflectCurrentUser := reflect.ValueOf(currentUser)
	reflectReqBody := reflect.ValueOf(reqBody)

	for i := 0; i < reflectReqBody.NumField(); i++ {
		if !reflectReqBody.Field(i).IsNil() {
			reflectCurrentUser.FieldByName(reflectReqBody.Type().Field(i).Name).Set(
				reflectReqBody.Field(i),
			)
		}
	}

	//logs.LogError(json.Unmarshal(reqBodyJSON, &currentUser))

	logs.AsJSON(reflectCurrentUser)

	//updRes, err := usersColl.ReplaceOne(context.Background(), bson.D{{"_id", oid}}, currentUser)
	logs.LogError(err)

	//fmt.Println(updRes)
}
