package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func saveLogFile(path string, data []byte) {
	var strData any
	logs.LogError(json.Unmarshal(data, &strData))
	writeData, err := json.MarshalIndent(strData, "", "  ")
	logs.LogError(err)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer func() {
		logs.LogError(f.Close())
	}()
	logs.LogError(f.Truncate(0))
	_, err = f.Write(writeData)
	logs.LogError(err)
}

func TestGetUserList(t *testing.T) {
	t.Run("user list 403", func(t *testing.T) {
		res, err := http.Get("http://localhost:8080/users")
		logs.LogError(err)

		body, err := io.ReadAll(res.Body)
		logs.LogError(err)
		if len(body) != 0 {
			logFilePath := "logs/test-log-getUsersLIST-" + strconv.Itoa(0) + ".json"
			saveLogFile(logFilePath, body)
		}

		if res.StatusCode != http.StatusForbidden {
			t.Error("code should be 403, but [" + strconv.Itoa(res.StatusCode) + "] got")
		}
	})
	t.Run("user list 200", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080/users", nil)
		req.Header.Add("X-admin-access", os.Getenv("ADMIN_ACCESS"))
		logs.LogError(err)
		res, err := http.DefaultClient.Do(req)
		logs.LogError(err)

		body, err := io.ReadAll(res.Body)
		logs.LogError(err)

		if len(body) != 0 {
			logFilePath := "logs/test-log-getUsersLIST-" + strconv.Itoa(1) + ".json"
			saveLogFile(logFilePath, body)
		}

		if res.StatusCode != http.StatusOK {
			t.Error("code should be 200, but [" + strconv.Itoa(res.StatusCode) + "] got")
		}
	})
	// 403
	// 200
}

func TestUserData(t *testing.T) {
	ids := []string{"6412de0ade9d0225a63fb0f4", "6412de0ade9d0225a63fb0f5"}

	for i, id := range ids {
		t.Run(id, func(t *testing.T) {
			res, err := http.Get("http://localhost:8080/users/" + id)
			logs.LogError(err)
			if res.StatusCode != 200 {
				t.Error("code is not 200, it is --- [" + strconv.Itoa(res.StatusCode) + "]")
			}
			body, err := io.ReadAll(res.Body)
			logs.LogError(err)

			logFilePath := "logs/test-log-getUserData-" + strconv.Itoa(i) + ".json"
			saveLogFile(logFilePath, body)
		})
	}
}

type postUserS struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO: replace InsertOne to InsertMany for postUser/2
// TODO: add description of the test to testData structure
func TestPostUser(t *testing.T) {

	testData := map[postUserS]int{
		postUserS{ // wrong data
			Name:     "",
			Email:    "",
			Password: "",
		}: 400,
		{ // wrong data
			Name:     "123",
			Email:    "123",
			Password: "123",
		}: 400,
		{ // user Exists
			Name:     "123",
			Email:    "123@gmail.com",
			Password: "123",
		}: 400,
	}
	rand.Seed(time.Now().UnixNano())
	successTestName := RandString(10, 3)
	testData[postUserS{ // Success
		Name:     successTestName,
		Email:    successTestName + "@gmail.com",
		Password: successTestName,
	}] = 200

	cnt := 0
	for test, waited := range testData {
		t.Run(strconv.Itoa(cnt), func(t *testing.T) {
			testJSON, err := json.Marshal(test)
			logs.LogError(err)

			res, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewReader(testJSON))
			logs.LogError(err)
			if res.StatusCode != waited {
				t.Error("code is not " + strconv.Itoa(waited) + ", it is --- [" + strconv.Itoa(res.StatusCode) + "]")
			}

			body, err := io.ReadAll(res.Body)
			logs.LogError(err)

			if len(body) != 0 {
				logFilePath := "logs/test-log-postUser-" + strconv.Itoa(cnt) + ".json"
				saveLogFile(logFilePath, body)
			}

			if waited == http.StatusOK {
				var inserted mongo.InsertOneResult
				logs.LogError(json.Unmarshal(body, &inserted))
				oid, err := primitive.ObjectIDFromHex(inserted.InsertedID.(string))
				logs.LogError(err)
				_, err = client.Database("todos").Collection("users").DeleteOne(context.TODO(), bson.D{{"_id", oid}})
				logs.LogError(err)
			}
		})
		cnt++
	}
}
