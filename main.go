package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"reflect"
)

//var QUERY_CONTEXT, CANCEL = context.WithTimeout(context.Background(), time.Second*10)

var client *mongo.Client

func Select(from *mongo.Collection, ctx context.Context, filter bson.D, res any) error {
	if reflect.TypeOf(res).Kind() != reflect.Ptr {
		fmt.Println(reflect.ValueOf(res).Kind())
		return errors.New("`res` must be pointer")
	}
	data, err := from.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return err
	}

	err = data.All(context.TODO(), res)
	return err
}

func init() {
	logs.LogError(godotenv.Load(".env"))
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")

	uri := "mongodb+srv://" + DbUser + ":" + DbPassword + "@test.acviawj.mongodb.net/?retryWrites=true&w=majority"
	err := errors.New("")
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	logs.LogError(err)
	logs.LogError(client.Ping(context.TODO(), nil))
	logs.LogSuccess("Connected to database\n")
}

func main() {
	defer func() {
		logs.LogError(client.Disconnect(context.TODO()))
		fmt.Println()
		logs.LogSuccess("Connection closed\n")
	}()

	router := gin.Default()

	router.GET("/users", getUserList)
	router.GET("/users/:id/data", getUserData)
	router.GET("/users/:id/info", getUserInfo)

	router.POST("users", postUser)
	// user image

	//todos

	logs.LogError(router.Run("localhost:8080"))

	logs.LogSuccess("SERVER started on `localhost:8080`")

	//initRandomData(9, 14, 22, true)
}
