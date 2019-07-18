# GoJAPI

GoJAPI(Go JSON API) is a small framework aimed at making a low level easy to use JSON API. It aims to do this by always returning JSON in the same basic format: Status, Data and Error. 

GoJAPI currently supports go1.10.2.
GoJAPI currently only supports MYSQL database

## Quick start

Navigate to Gopath/src/github.com:

* On Windows: `cd %gopath%/src/github.com`

* On Linux: `cd $GOPATH/src/github.com`


If not already exist create a folder here where you want to work on the project:
```bash
mkdir gojapi
cd gojapi
```


Clone this repository to your gopath:
```bash
git clone https://github.com/IanVinkHub/gojapi.git
```


Now goto the cloned project:
```bash
cd gojapi
```

> Before continuing follow the Configuration section under this and return

Build and run:
* On Windows:
`go build & gojapi`

* On Linux:
`go build & ./gojapi`


Now if you go to <http://localhost:8080> you shall see the following:
```JSON
{"status":"error","data":"","error":{"code":404,"custom_code":-1,"type":"Not Found","param":[],"message":"Endpoint not found"}}
```

## Configuration

### Database

First we need to setup the database by default the framework will use the gojapi database with root on port 3306 without password this can be configured in `controllers/init.go`.

To change this you can change: 
```Go
	qb.DBC = db.Connection{Database: "gojapi"}
	qb.DBC.Default()
```
to
```Go
	qb.DBC = db.Connection{
		Host: "HOST",
		Port: "PORT",
		Database: "DATABASE",
		Username: "USERNAME",
		Password: "PASSWORD",
		}
```
If you don't want to change one of them just leave them out, only database is required.


## File Structure

In all these folders you can find a `README.md` with extra details

	-app -> Here are all general functions stored
	-controllers -> Here are all serve functions stored
	-http -> Here are all examples for API Calls stored
	-models -> Here are all models stored for the db and requests/responses
	-sql -> Here are all sql files stored to quickly setup the database

## Code Preview

This framework does not add a different webserver it just adds extra functionality for using JSON and connecting to the database.

---

### Responding with data

In this framework most data is written to the client using jio.WriteMsg example:
```GO
/*
type LoginPost struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Auth     TokenPost `json:"auth"`
}
type TokenPost struct {
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	ExpireAt  string `json:"expire_at"`
	DeleteAt  string `json:"delete_at"`
}
*/
// Here returnedRow is the LoginPost struct above filled
jio.WriteMsg(w, returnedRow, "")
```
Returns
```JSON
{
    "status": "succes",
    "data": {
        "id": 1,
        "username": "Legend27",
        "auth": {
            "token": "rQu)I)J3yVJn9v:SBj3R0YiynQvRN4-D6h0!Yxd)dXXt!Vr(X6+ZUmzEb(qfnU)S",
            "created_at": "Wed, 17 Jul 2019 22:07:10 +1700",
            "expire_at": "Wed, 17 Jul 2019 22:07:50 +1700",
            "delete_at": "Wed, 17 Jul 2019 22:08:10 +1700"
        }
    },
    "error": ""
}
```

---

### Checking parameters

This framework comes with a function to check for required parameters, it will return a list of missing parameters:

#### Example where it checks if all properties are filled
```Go
// Check if any params are invalid
errorParam := jio.CheckParamsStruc(body)
if len(errorParam) > 0 {
	errMsg := jio.CreateErrorMsg(jio.InvalidParameters, errorParam)
	jio.WriteMsg(w, "", errMsg)
	return
}
```
#### Example where you select wich are required:
```Go
// Check if any required params are invalid
errorParam := jio.CheckParams(body, [1]String{"Username"})
if len(errorParam) > 0 {
	errMsg := jio.CreateErrorMsg(jio.InvalidParameters, errorParam)
	jio.WriteMsg(w, "", errMsg)
	return
}
```

---

### Responding with Errors

In this framework we add new functions for responding with errors here are a few examples:


#### Using jio.Error -> Used for common errors that don't need more explanation then default.
```Go
// methodNotAllowedMessage = "Method not allowed for this API endpoint"
jio.Error(w, errors.New(methodNotAllowedMessage), 405)
```
Returns
```JSON
{
	"status":"error",
	"data":"",
	"error":{
		"code":405,
		"custom_code":-1,
		"type":"Method Not Allowed",
		"param":[],
		"message":"Method not allowed for this API endpoint"
	}
}
```

#### Using jio.CreateErrorMsg -> This is ussualy used for custom errors
```Go
/*
 You can edit and add custom errors in app/jio/out.go read app/jio/README.md for more info
*/
// jio.TokenExpired = 4001
errMsg := jio.CreateErrorMsg(jio.InvalidParameters, []string{})
jio.WriteMsg(w, "", errMsg)
```
Returns (body = {})
```JSON
{
    "status": "error",
    "data": "",
    "error": {
        "code": 400,
        "custom_code": 4001,
        "type": "Invalid parameters",
        "param": [
            "username",
            "password"
        ],
        "message": "The following parameters does not exist or empty: [username, password]"
    }
}
```

#### using jio.CreateErrorMsgError -> this is ussualy used for custom error and with a Error type
```Go
/*
 You can edit and add custom errors in app/jio/out.go read app/jio/README.md for more info
*/
// err = json unmarshal error
// jio.InvalidJSON = 4002
errMsg := jio.CreateErrorMsgError(jio.InvalidJSON, []string{}, err)
jio.WriteMsg(w, "", errMsg)
```
Returns (body = )
```JSON
{
    "status": "error",
    "data": "",
    "error": {
        "code": 400,
        "custom_code": 4002,
        "type": "Invalid JSON",
        "param": [],
        "message": "unexpected end of JSON input"
    }
}
```

---

### Getting structs from database 
This framework adds a function wich makes it easier to convert sql.rows to any object you want.

#### Example where all users are selected
```go
// Get all users
rows, err := qb.From("users").WhereValue("username", body.Username).All()
if err != nil {
	jio.Error(w, err, 500)
	return nil
}

userRows := db.ReadRows(dbmodels.User{}, rows)
```
### Querying the database

In this framework there are functions added wich make it easier to query the database.

#### Example getting users
```GO
// This will get all users where username equals body.username
/*
Note that the second parameter in wherevalue always gets escaped with html escapestring to prevent SQL Injections
*/
rows, err := qb.From("users").WhereValue("username", body.Username).All()
if err != nil {
	jio.Error(w, err, 500)
	return nil
}
```

### Example inserting user
```GO
// Notice that you only need to give the table and a struct
rows, err = qb.Insert("users", body)
if err != nil {
	jio.Error(w, err, 500)
	return
}
```

### Example update user
```GO
_, err = qb.WhereValue("id", id).Update("users", body)
if err != nil {
	jio.Error(w, err, 500)
	return
}
```

---

## Code Documentation

The full code documentation has not been added yet I highly recommend just checking trough the code since it is not a lot, and you will get a better view of what parts you need.
