# todoAPI

[![CodeFactor](https://www.codefactor.io/repository/github/readyyyk/todoapi/badge)](https://www.codefactor.io/repository/github/readyyyk/todoapi)
[![Go report](https://goreportcard.com/badge/github.com/readyyyk/todoAPI)](https://goreportcard.com/report/github.com/readyyyk/todoAPI)

### Custom error descriptions

| Code | Meaning                             |
|------|-------------------------------------|
| 0    | invalid data                        |
| 1    | User with this email already exists |
| 2    | user don't exists                   |
| 3    | Wrong password                      |
| 4    | JWT token is invalid                |
| 5    | Group don't exist                   |
| 6    | User doesn't own this group         |

### Routes

| Action | Method | Route | Req. data |
|----|:---:| ---- | ---- |
| _**USER**_ |
| Create user| POST | <host>/users | `body` {"email" string, "name" string} |
| Get user list | GET | <host>/users | `header` "X-admin-access" |
| Get user info | GET | <host>/users/:id | `none` |
| Get user data | GET | <host>/users/:id/data | `none` |
| Update user data | POST | <host>/users/:id | `body` {"email" string, "password" base64, "name" string } |
| Delete user | DELETE | <host>/users/:id | `none` |
|  Login | POST  | <host>/users/login | `body` {"email" string, "password" base64} |
|  _**GROUP**_ |
| Create group | POST | <host>/groups | `body` {"title" string} |
| Delete group | DELETE | <host>/groups/:id | `JWT` |
| _**TODO**_  |
| Create todo | POST | <host>/groups/:group_id | `body` {"title" string, "text" string, "deadline" date} |
| Delete todo | DELETE | <host>/groups/:group_id/todos/:todo_id | `JWT` |
