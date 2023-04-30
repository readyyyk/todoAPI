package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"github.com/readyyyk/todoAPI/m_group"
	"github.com/readyyyk/todoAPI/m_todo"
	"github.com/readyyyk/todoAPI/m_user"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

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

func init() {
	logs.LogError(godotenv.Load(".env"))
}

func main() {
	logs.LogSuccess(proceeding.GenerateJWT(primitive.NewObjectID()) + "\n")

	// Defining routes
	router := gin.Default()
	router.Use(cors.Default())

	// ! API routes
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, "github.com/readyyyk/todoAPI") })

		// TODO: user image
		// ! USERS
		usersRoutes := apiRoutes.Group("/users")
		{
			usersRoutes.POST(routes.user.c, m_user.Create)       // c
			usersRoutes.GET(routes.user.getList, m_user.GetList) // get list of all users
			usersRoutes.GET(routes.user.r, m_user.GetInfo)       // r
			usersRoutes.GET(routes.user.getData, m_user.GetData) // get entire data
			usersRoutes.PUT(routes.user.u, m_user.Update)        // u
			usersRoutes.DELETE(routes.user.d, m_user.Delete)     // d
			usersRoutes.POST(routes.user.getLogin, m_user.Login) // auth
		}

		// ! GROUPS	access only for owners
		//			"Auth" header required
		groupRoutes := apiRoutes.Group("/groups")
		{
			groupRoutes.POST(routes.groups.c, m_group.Create)
			groupRoutes.DELETE(routes.groups.d, m_group.Delete)
			// router.PUT("/:id", updateGroup)
			// router.GET("/:id", getGroup)

			// ! TODOS
			todoRoutes := groupRoutes.Group("/:group_id/todos")
			{
				todoRoutes.POST(routes.todos.c, m_todo.Create)
				todoRoutes.DELETE(routes.todos.d, m_todo.Delete)
				//todoRoutes.GET("/:group_id/todos/:todo_id", getTodo)
				//todoRoutes.PUT("/:group_id/todos/:todo_id", updateTodo)
			}
		}
	}

	//initRandomData(9, 14, 22, true)
	logs.LogSuccess("Database url: " + os.Getenv("DB_URL") + "\n")
	logs.LogSuccess("SERVER starting on `" + host + "`...\n")
	logs.LogError(router.Run(host))
}
