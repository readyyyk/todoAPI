package proceeding

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

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
