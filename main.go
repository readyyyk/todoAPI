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
	"net/http"
	"os"
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

// <host>/api/...
var routes = routesT{
	// <host>/api/users/...
	user: userRoutesT{
		routesCRUD: routesCRUD{
			c: "",
			r: "/:id",
			u: "/:id",
			d: "/:id",
		},
		getLogin: "/login",
		getData:  "/:id/data",
		getList:  "",
	},

	// Todo: groups - Add ["add user", "rec deletion", "update"], unit tests
	// <host>/api/groups/...
	groups: routesCRUD{
		c: "",
		d: "/:group_id",
	},

	// Todo: groups - Add ["update"], unit tests
	// <host>/api/groups/:group_id/todos/...
	todos: routesCRUD{
		c: "",
		d: "/:todo_id",
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

	AvailableTodoState = []string{"passive", "ongoing", "done", "important", "expired"}

	url := fmt.Sprintf(os.Getenv("DB_URL")) //, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))

	err := errors.New("")
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	logs.LogError(err)
	logs.LogError(client.Ping(context.TODO(), nil))
	logs.LogSuccess("Connected to database\n")
}

func main() {
	defer func() {
		logs.LogError(client.Disconnect(context.TODO()))
		logs.LogSuccess("\nConnection closed\n")
	}()

	logs.LogSuccess(GenerateJWT(primitive.NewObjectID()))

	// Defining routes
	router := gin.Default()

	// ! API routes
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "github.com/readyyyk/todoAPI")
		})

		// TODO: user image
		// ! USERS
		usersRoutes := apiRoutes.Group("/users")
		{
			usersRoutes.POST(routes.user.c, createUser)       // c
			usersRoutes.GET(routes.user.getList, getUserList) // get list of all users
			usersRoutes.GET(routes.user.r, getUserInfo)       // r
			usersRoutes.GET(routes.user.getData, getUserData) // get entire data
			usersRoutes.PUT(routes.user.u, updateUser)        // u
			usersRoutes.DELETE(routes.user.d, deleteUser)     // d
			usersRoutes.POST(routes.user.getLogin, loginUser) // auth
		}

		// ! GROUPS	access only for owners
		//			"Auth" header required
		groupRoutes := apiRoutes.Group("/groups")
		{
			groupRoutes.POST("/", createGroup)
			groupRoutes.DELETE("/:group_id", deleteGroup)
			// router.PUT("/:id", updateGroup)
			// router.GET("/:id", getGroup)

			// ! TODOS
			todoRoutes := groupRoutes.Group("/:group_id/todos")
			{
				todoRoutes.POST("/", createTodo)
				todoRoutes.DELETE("/:todo_id", deleteTodo)
				//todoRoutes.GET("/:group_id/todos/:todo_id", getTodo)
				//todoRoutes.PUT("/:group_id/todos/:todo_id", updateTodo)
			}
		}
	}

	//initRandomData(9, 14, 22, true)
	logs.LogSuccess("Connected database url: " + os.Getenv("DB_URL") + "\n")
	logs.LogSuccess("SERVER starting on `" + host + "`...\n")
	logs.LogError(router.Run(host))
}
