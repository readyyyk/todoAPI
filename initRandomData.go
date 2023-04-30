package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"github.com/readyyyk/todoAPI/pkg/proceeding"
	"github.com/readyyyk/todoAPI/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"sort"
	"time"
)

func RandString(n int, min int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_."
	b := make([]byte, n+min)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func initRandomData(usersCnt int, groupsCnt int, todosCnt int, logRes bool) {
	rand.Seed(time.Now().UnixNano())
	client := proceeding.NewDbClient()

	users := client.Database("todos").Collection("users")
	groups := client.Database("todos").Collection("groups")
	todos := client.Database("todos").Collection("todos")

	var userIds []primitive.ObjectID
	var groupIds []primitive.ObjectID

	var userDataSet []interface{}
	for i := 0; i < usersCnt; i++ {
		nID := primitive.NewObjectID()
		userIds = append(userIds, nID)

		pswd := base64.StdEncoding.EncodeToString([]byte(RandString(12, 8)))

		newUser := types.User{
			Id:         nID,
			Email:      RandString(rand.Intn(20), 2) + "@gmail.com",
			Password:   pswd,
			Registered: time.Now(),
			Name:       RandString(10, 1),
		}
		userDataSet = append(userDataSet, newUser)
	}
	//logs.AsJSON(userDataSet)

	var groupDataSet []interface{}
	for i := 0; i < groupsCnt; i++ {
		nID := primitive.NewObjectID()
		groupIds = append(groupIds, nID)

		// make a copy of userIds to avoid changing all elems in GroupDataSet to the state of last
		userIdsCopy := make([]primitive.ObjectID, len(userIds))
		copy(userIdsCopy, userIds)

		rand.Shuffle(len(userIdsCopy), func(i, j int) {
			userIdsCopy[i], userIdsCopy[j] = userIdsCopy[j], userIdsCopy[i]
		})
		//logs.AsJSON(userIds)

		ownersTmp := userIdsCopy[:rand.Intn(4)+1]
		sort.Slice(ownersTmp, func(i, j int) bool {
			return ownersTmp[i].String() < ownersTmp[j].String()
		})

		newGroup := types.Group{
			Id:      nID,
			Owners:  ownersTmp,
			Title:   RandString(30, 2),
			Created: time.Now(),
		}
		//logs.AsJSON(groupDataSet)

		groupDataSet = append(groupDataSet, newGroup)
	}
	//logs.AsJSON(groupDataSet)

	var todosDataSet []interface{}
	for i := 0; i < todosCnt; i++ {
		rand.Shuffle(len(groupIds), func(i, j int) {
			groupIds[i], groupIds[j] = groupIds[j], groupIds[i]
		})

		newTodo := types.Todo{
			Id:        primitive.NewObjectID(),
			Group:     groupIds[0],
			Title:     RandString(20, 1),
			Text:      RandString(100, 20),
			State:     "passive",
			StartDate: time.Now(),
			Deadline:  time.Now().Add(time.Hour * time.Duration(rand.Intn(300)+1)),
		}
		todosDataSet = append(todosDataSet, newTodo)
	}
	//logs.AsJSON(todosDataSet)

	//err := Client.Database("todos").Drop(context.TODO())
	//logs.LogError(err)

	logs.LogError(users.Drop(context.TODO()))
	logs.LogError(groups.Drop(context.TODO()))
	logs.LogError(todos.Drop(context.TODO()))

	resUsers, err := users.InsertMany(context.TODO(), userDataSet)
	logs.LogError(err)

	resGroups, err := groups.InsertMany(context.TODO(), groupDataSet)
	logs.LogError(err)

	resTodos, err := todos.InsertMany(context.TODO(), todosDataSet)
	logs.LogError(err)

	if logRes {
		resUsersJSON, err := json.MarshalIndent(resUsers, "", "  ")
		logs.LogError(err)
		fmt.Println(string(resUsersJSON))
		resGroupsJSON, err := json.MarshalIndent(resGroups, "", "  ")
		logs.LogError(err)
		fmt.Println(string(resGroupsJSON))
		resTodosJSON, err := json.MarshalIndent(resTodos, "", "  ")
		logs.LogError(err)
		fmt.Println(string(resTodosJSON))
	}
}
