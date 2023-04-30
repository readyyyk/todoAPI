# todoAPI

[![CodeFactor](https://www.codefactor.io/repository/github/readyyyk/todoapi/badge)](https://www.codefactor.io/repository/github/readyyyk/todoapi)
[![Go report](https://goreportcard.com/badge/github.com/readyyyk/todoAPI)](https://goreportcard.com/report/github.com/readyyyk/todoAPI)

### Custom error descriptions

| Code | Meaning                             |
|------|-------------------------------------|
| 0    | Invalid data                        |
| 1    | User with this email already exists |
| 2    | User don't exists                   |
| 3    | Wrong password                      |
| 4    | JWT token is invalid                |
| 5    | Group don't exist                   |
| 6    | User doesn't own this group         |

> ---
 
### Routes

| Action           | Method | Route                                       |
|------------------|:------:|---------------------------------------------|
| _**USER**_       |
| Create user      |  POST  | <host\>/api/users                           |
| Get user list    |  GET   | <host\>/api/users                           |
| Get user info    |  GET   | <host\>/api/users/:id                       |
| Get user data    |  GET   | <host\>/api/users/:id/data                  |
| Update user data |  POST  | <host\>/api/users/:id                       |
| Delete user      | DELETE | <host\>/api/users/:id                       |
| Login            |  POST  | <host\>/api/users/login                     |
| _**GROUP**_      |
| Create group     |  POST  | <host\>/api/groups                          |
| Delete group     | DELETE | <host\>/api/groups/:id                      |
| _**TODO**_       |
| Create todo      |  POST  | <host\>/api/groups/:group_id                |
| Delete todo      | DELETE | <host\>/api/groups/:group_id/todos/:todo_id |

> ---

### Provided data-sets

| Action           | _**Req.**_ data                                                    |
|------------------|--------------------------------------------------------------------|
| _**USER**_       |
| Create user      | `body` {"email" string, "name" string}                             |
| Get user list    | `header` "X-admin-access"                                          |
| Get user info    | `none`                                                             |
| Get user data    | `JWT`                                                              |
| Update user data | `JWT`, `body` {"email" string, "password" base64, "name" string }  |
| Delete user      | `JWT`                                                              |
| Login            | `body` {"email" string, "password" base64}                         |
| _**GROUP**_      |
| Create group     | `JWT`, `body` {"title" string}                                     |
| Delete group     | `JWT`                                                              |
| _**TODO**_       |
| Create todo      | `JWT`, `body` {"title" string, "text" string, "deadline" date}     |
| Delete todo      | `JWT`                                                              |

> ---
> 🚨 __*Will be added during writing tests*__

| Action           | _**Resp.**_ data |
|------------------|------------------|
| _**USER**_       |                  |
| Create user      |                  |
| Get user list    |                  |
| Get user info    |                  |
| Get user data    |                  |
| Update user data |                  |
| Delete user      |                  |
| Login            |                  |
| _**GROUP**_      |                  |
| Create group     |                  |
| Delete group     |                  |
| _**TODO**_       |                  |
| Create todo      |                  |
| Delete todo      |                  |

> ---
 
### Files
    
| Action           | File                                                                           | State                |
|------------------|--------------------------------------------------------------------------------|----------------------|
| _**USER**_       |
| Create user      | [/user_post](https://github.com/readyyyk/todoAPI/blob/main/user_post.go)       | 🚨 Need to be tested |
| Get user list    | [/user_getList](https://github.com/readyyyk/todoAPI/blob/main/user_getList.go) | 🦺 Need to be tested |
| Get user info    | [/user_getInfo](https://github.com/readyyyk/todoAPI/blob/main/user_getList.go) | 🦺 Need to be tested |
| Get user data    | [/user_getData](https://github.com/readyyyk/todoAPI/blob/main/user_getData.go) | 🚨 Need to be tested |
| Update user data | [/user_upd](https://github.com/readyyyk/todoAPI/blob/main/user_upd.go)         | 🦺 Need to be tested |
| Delete user      | [/user_del](https://github.com/readyyyk/todoAPI/blob/main/user_del.go)         | 🦺 Need to be tested |
| Login            | [/user_login](https://github.com/readyyyk/todoAPI/blob/main/user_login.go)     | 🚨 Need to be tested |
| _**GROUP**_      |
| Create group     | [/group_post](https://github.com/readyyyk/todoAPI/blob/main/group_post.go)     | 🚨 Need to be tested | 
| Delete group     | [/group_delete](https://github.com/readyyyk/todoAPI/blob/main/group_delete.go) | 🚨 Need to be tested | 
| _**TODO**_       |
| Create todo      | [/todo_post](https://github.com/readyyyk/todoAPI/blob/main/todo_post.go)       | 🚨 Need to be tested |
| Delete todo      | [/todo_delete](https://github.com/readyyyk/todoAPI/blob/main/todo_delete.go)   | 🚨 Need to be tested |