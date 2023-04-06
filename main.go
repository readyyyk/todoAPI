package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

var client *mongo.Client

const host = "localhost:8080"

type routesCRUD struct {
	c, r, u, d string
}
type userRoutesT struct {
	routesCRUD
	getData  string
	getList  string
	getLogin string
}
type routesT struct {
	user   userRoutesT
	groups routesCRUD
	todos  routesCRUD
}

var routes = routesT{
	user: userRoutesT{
		routesCRUD: routesCRUD{
			c: "/users",
			r: "/users/:id",
			u: "/users/:id",
			d: "/users/:id",
		},
		getLogin: "/users/:email/:password",
		getData:  "/users/login",
		getList:  "/users",
	},
}

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
	//DbUser := os.Getenv("DB_USER")
	//DbPassword := os.Getenv("DB_PASSWORD")

	//url := "mongodb+srv://" + DbUser + ":" + DbPassword + "@test.acviawj.mongodb.net/?retryWrites=true&w=majority"
	url := "mongodb://127.0.0.1:27017"

	err := errors.New("")
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
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

	logs.LogSuccess(GenerateJWT(primitive.NewObjectID()))
	router := gin.Default()

	// user
	// TODO user image
	router.POST(routes.user.c, postUser)         // c
	router.GET(routes.user.getList, getUserList) // get list of all users
	router.GET(routes.user.r, getUserInfo)       // r
	router.GET(routes.user.getData, getUserData) // get entire data
	router.PUT(routes.user.u, updateUser)        // u
	router.DELETE(routes.user.d, deleteUser)     // d
	router.POST(routes.user.getData, loginUser)  // auth
	/*
		// groups	access only for owner
		//			Logged header required
		router.POST("/groups", createGroup)
		router.GET("/groups/:id", getGroup)
		router.PUT("/groups/:id", updateGroup)
		router.DELETE("/groups/:id", deleteGroup)
		// PUT /groups/:id/users/:userId

		// todos
		router.POST("/groups/:group_id/todos", createTodo)
		router.GET("/groups/:group_id/todos/:todo_id", getTodo)
		router.PUT("/groups/:group_id/todos/:todo_id", updateTodo)
		router.DELETE("/groups/:group_id/todos/:todo_id", deleteTodo)
	*/

	//initRandomData(9, 14, 22, true)

	logs.LogWarning("SERVER starting on `" + host + "`...\n")
	logs.LogError(router.Run(host))

}
