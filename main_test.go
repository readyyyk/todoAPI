package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/readyyyk/terminal-todos-go/pkg/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
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
			logFilePath := "logs/test-log-getUsersLIST--" + strconv.Itoa(0) + ".json"
			saveLogFile(logFilePath, body)
		}

		if res.StatusCode != http.StatusForbidden {
			t.Error("Code should be 403, but [" + strconv.Itoa(res.StatusCode) + "] got")
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
			logFilePath := "logs/test-log-getUsersLIST--" + strconv.Itoa(1) + ".json"
			saveLogFile(logFilePath, body)
		}

		if res.StatusCode != http.StatusOK {
			t.Error("Code should be 200, but [" + strconv.Itoa(res.StatusCode) + "] got")
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
				t.Error("Code is not 200, it is --- [" + strconv.Itoa(res.StatusCode) + "]")
			}
			body, err := io.ReadAll(res.Body)
			logs.LogError(err)

			logFilePath := "logs/test-log-getUserData--" + strconv.Itoa(i) + ".json"
			saveLogFile(logFilePath, body)
		})
	}
}

type testDataT struct {
	data        []testUserPostDataT
	resp        []testDataRespWaitedT
	description string
}
type testUserPostDataT struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type testDataRespWaitedT struct {
	statusCode int
	err        errorDescriptionT
}

//	type errorDescriptionT struct {
//		Code        int    `json:"code"`
//		Description string `json:"description"`
//	}
func TestPostUser(t *testing.T) {

	/*
		{
			data 					testUserPostData
			resp {
				statusCode 			int
				error
					Code			int
		 			Description 	string
			}
		}
	*/
	rand.Seed(time.Now().UnixNano())
	successTestName := RandString(10, 3)

	testData := []testDataT{
		// wrong data
		{
			data: []testUserPostDataT{
				{
					Name:     "",
					Email:    "",
					Password: "",
				},
			},
			resp: []testDataRespWaitedT{
				{
					statusCode: 400,
					err: errorDescriptionT{
						Code:        0,
						Description: "Invalid data",
					},
				},
			},
			description: "wrong data",
		},
		{
			data: []testUserPostDataT{
				{
					Name:     "123",
					Email:    "123",
					Password: "123",
				},
			},
			resp: []testDataRespWaitedT{
				{
					statusCode: 400,
					err: errorDescriptionT{
						Code:        0,
						Description: "Invalid data",
					},
				},
			},
			description: "wrong data",
		},

		// user Exists
		{
			data: []testUserPostDataT{
				{
					Name:     successTestName,
					Email:    successTestName + "@gmail.com",
					Password: successTestName,
				},
				{
					Name:     successTestName,
					Email:    successTestName + "@gmail.com",
					Password: successTestName,
				},
			},
			resp: []testDataRespWaitedT{
				{
					statusCode: 200,
					err: errorDescriptionT{
						Code:        -1,
						Description: "",
					},
				},
				{
					statusCode: 400,
					err: errorDescriptionT{
						Code:        1,
						Description: "User with this email already exists",
					},
				},
			},
			description: "existing user",
		},

		// success
		{
			data: []testUserPostDataT{
				{
					Name:     successTestName,
					Email:    successTestName + "@gmail.com",
					Password: successTestName,
				},
			},
			resp: []testDataRespWaitedT{
				{
					statusCode: 200,
					err: errorDescriptionT{
						Code:        -1,
						Description: "",
					},
				},
			},
			description: "success",
		},
	}

	for testNumber, test := range testData {
		t.Run(test.description, func(t *testing.T) {

			var toDeleteIds []primitive.ObjectID

			for i, currentTest := range test.data {
				waited := test.resp[i]

				currentTestJSON, err := json.Marshal(currentTest)
				logs.LogError(err)

				res, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewReader(currentTestJSON))
				logs.LogError(err)
				if res.StatusCode != waited.statusCode {
					t.Error("Code is not " + strconv.Itoa(waited.err.Code) + ", it is --- [" + strconv.Itoa(res.StatusCode) + "]")
				}

				body, err := io.ReadAll(res.Body)
				logs.LogError(err)

				var responseData any
				err = json.Unmarshal(body, &responseData)
				logs.LogError(err)
				//logs.Deb(reflect.TypeOf(responseData).String())

				if responseData.(map[string]interface{})["InsertedID"] != nil {
					oid, err := primitive.ObjectIDFromHex(responseData.(map[string]interface{})["InsertedID"].(string))
					logs.LogError(err)
					toDeleteIds = append(toDeleteIds, oid)
				} else {
					//fmt.Println(responseData)
					var waitedErrMap map[string]interface{}
					waitedErrMapJSON, _ := json.Marshal(waited.err)
					_ = json.Unmarshal(waitedErrMapJSON, &waitedErrMap)

					if !reflect.DeepEqual(responseData, waitedErrMap) {
						t.Error("DeepEqual returned false:\n", responseData)
					}
				}

				if len(body) != 0 {
					logFilePath := "logs/test-postUser--" + strconv.Itoa(testNumber) + "-" + strconv.Itoa(i) + ".json"
					saveLogFile(logFilePath, body)
				}
			}

			for _, oid := range toDeleteIds {
				_, err := client.Database("todos").Collection("users").DeleteOne(context.TODO(), bson.D{{"_id", oid}})
				logs.LogError(err)
			}
		})
	}
}
