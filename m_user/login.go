package m_user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	apiErrors "github.com/readyyyk/todoAPI/pkg/errors"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"github.com/readyyyk/todoAPI/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	type loginRespT struct {
		Logged bool                        `json:"logged"`
		Err    apiErrors.ErrorDescriptionT `json:"err"`
		Token  string                      `json:"token"`
	}
	type userLoginT struct {
		Email string `json:"email" validate:"required,email"`
		Pswd  string `json:"password" validate:"required,base64"`
	}

	var userData userLoginT
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		return
	}

	err = json.Unmarshal(jsonData, &userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiErrors.Errors[0])
		return
	}

	uemail := userData.Email
	upswd, err := base64.StdEncoding.DecodeString(userData.Pswd)

	// validate entered data
	if _, ok := err.(base64.CorruptInputError); ok || len(uemail) == 0 || len(upswd) == 0 {
		c.JSON(http.StatusBadRequest, loginRespT{
			false,
			apiErrors.Errors[0],
			"",
		})
		return
	} /*else if len(uemail) == 0 || len(upswd) == 0 {
		c.JSON(http.StatusBadRequest, errorDescriptionT{
			Code:        0,
			Description: "Cannot get username or password",
		})
		return
	}*/

	// select user with provided email
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	var userSelectRes []types.User
	client := proceeding.NewDbClient()
	logs.LogError(proceeding.Select(
		client.Database("todos").Collection("users"),
		ctx,
		bson.D{{"email", uemail}},
		&userSelectRes,
	))
	if len(userSelectRes) == 0 {
		c.JSON(http.StatusNotFound, loginRespT{
			false,
			apiErrors.Errors[2],
			"",
		})
		return
	}
	userSelected := userSelectRes[0]

	// validate password
	if userSelected.Password != base64.StdEncoding.EncodeToString(upswd) {
		c.JSON(http.StatusNotFound, loginRespT{
			false,
			apiErrors.Errors[3],
			"",
		})
		return
	}

	// generate jwt
	signedToken := proceeding.GenerateJWT(userSelected.Id)

	c.JSON(http.StatusOK, loginRespT{
		true,
		apiErrors.ErrorDescriptionT{},
		signedToken,
	})
}
