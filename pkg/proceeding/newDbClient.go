package proceeding

import (
	"context"
	"fmt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func NewDbClient() *mongo.Client {
	url := fmt.Sprintf(os.Getenv("DB_URL")) //, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	logs.LogError(err)
	logs.LogError(client.Ping(context.TODO(), nil))
	return client
}
