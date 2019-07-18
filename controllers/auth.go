package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/IanVinkHub/gojapi/app"
	"github.com/IanVinkHub/gojapi/app/db"
	"github.com/IanVinkHub/gojapi/app/jio"
	"github.com/IanVinkHub/gojapi/models/dbmodels"
	"github.com/IanVinkHub/gojapi/models/reqmodels"
	"github.com/IanVinkHub/gojapi/models/resmodels"

	"golang.org/x/crypto/bcrypt"
)

const cost = 12
const methodNotAllowedMessage = "Method not allowed for this API endpoint"

func getUserLogin(w http.ResponseWriter, body reqmodels.LoginPost) *resmodels.LoginPost {
	// Get all users
	rows, err := qb.From("users").WhereValue("username", body.Username).All()
	if err != nil {
		jio.Error(w, err, 500)
		return nil
	}

	//returnedRows
	var returnedRow resmodels.LoginPost
	userRows := db.ReadRows(dbmodels.User{}, rows)

	if len(userRows) == 0 {
		// Comparing hash for adding realistic delay to try to prevent username checking for existance by delay
		bcrypt.CompareHashAndPassword([]byte("Password123456"), []byte(body.Password))

		errMsg := jio.CreateErrorMsg(jio.CredentialsNotFound, nil)
		jio.WriteMsg(w, "", errMsg)
		return nil
	}

	if len(userRows) > 1 {
		jio.Error(w, err, 500)
		return nil
	}

	user := userRows[0].(dbmodels.User)

	// Check if correct password was given
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		errMsg := jio.CreateErrorMsg(jio.CredentialsNotFound, nil)
		jio.WriteMsg(w, "", errMsg)
		return nil
	}

	token := app.CreateUserToken(user.Username)
	userToken, errMsg := app.GetUserToken(user.Username, token)
	if errMsg != nil {
		jio.WriteMsg(w, "", errMsg)
		return nil
	}

	returnedRow = resmodels.LoginPost{ID: user.ID, Username: user.Username, Auth: userToken.GetTokenPostResponse()}

	return &returnedRow
}

// Login handles the login api call
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jio.Error(w, errors.New(methodNotAllowedMessage), 405)
		return
	}

	// Read body
	b := jio.Readbody(w, r)
	if b == nil {
		return
	}

	// Unmarshal body to body variable
	var body reqmodels.LoginPost
	err := json.Unmarshal(b, &body)
	if err != nil {
		errMsg := jio.CreateErrorMsgError(jio.InvalidJSON, []string{}, err)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	// Check if any params are invalid
	errorParam := jio.CheckParamsStruc(body)
	if len(errorParam) > 0 {
		errMsg := jio.CreateErrorMsg(jio.InvalidParameters, errorParam)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	returnedRow := getUserLogin(w, body)
	if returnedRow == nil {
		return
	}

	fmt.Println(returnedRow.Auth.Token)
	jio.WriteMsg(w, returnedRow, "")
}

// Register handles the register api call
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jio.Error(w, errors.New(methodNotAllowedMessage), 405)
		return
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		jio.Error(w, err, 500)
		return
	}

	// Unmarshal body to body variable
	var body reqmodels.LoginPost
	err = json.Unmarshal(b, &body)
	if err != nil {
		errMsg := jio.CreateErrorMsgError(jio.InvalidJSON, []string{}, err)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	// Check if any params are invalid
	errorParam := jio.CheckParamsStruc(body)
	if len(errorParam) > 0 {
		errMsg := jio.CreateErrorMsg(jio.InvalidParameters, errorParam)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	// Get all users
	rows, err := qb.From("users").WhereValue("username", body.Username).All()
	if err != nil {
		jio.Error(w, err, 500)
		return
	}

	// Check if username not yet used
	userRows := db.ReadRows(dbmodels.User{}, rows)
	if len(userRows) != 0 {
		errMsg := jio.CreateErrorMsg(jio.CredentialsExist, nil)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	// Encrypt password
	psw, err := bcrypt.GenerateFromPassword([]byte(body.Password), cost)
	if err != nil {
		jio.Error(w, err, 500)
		return
	}
	oldpsw := body.Password
	body.Password = string(psw)

	// Insert user data
	rows, err = qb.Insert("users", body)
	if err != nil {
		jio.Error(w, err, 500)
		return
	}

	// Gets the data to return and makes sure password won't return
	fmt.Println(body.Username)
	fmt.Println(body.Password)
	body.Password = oldpsw
	returnedRow := getUserLogin(w, body)
	if returnedRow == nil {
		return
	}

	fmt.Println(returnedRow.Auth.Token)
	jio.WriteMsg(w, returnedRow, "")
}

// Changepassword handles the changepassword api call
func Changepassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jio.Error(w, errors.New(methodNotAllowedMessage), 405)
		return
	}

	// Read body
	b := jio.Readbody(w, r)
	if b == nil {
		return
	}

	// Unmarshal body to body variable
	var body reqmodels.ChangePasswordPost
	err := json.Unmarshal(b, &body)
	if err != nil {
		errMsg := jio.CreateErrorMsgError(jio.InvalidJSON, []string{}, err)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	// Check if any params are invalid
	errorParam := jio.CheckParamsStruc(body)
	if len(errorParam) > 0 {
		errMsg := jio.CreateErrorMsg(jio.InvalidParameters, errorParam)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	returnedRow := getUserLogin(w, reqmodels.LoginPost{Username: body.Username, Password: body.Password})
	if returnedRow == nil {
		return
	}

	// Encrypt password
	psw, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), cost)
	if err != nil {
		jio.Error(w, err, 500)
		return
	}

	_, err = qb.WhereValue("id", strconv.Itoa(returnedRow.ID)).Update("users", reqmodels.LoginPost{Username: body.Username, Password: string(psw)})
	if err != nil {
		jio.Error(w, err, 500)
		return
	}

	// Gets the data to return and makes sure password won't return
	/*
		returnRows := []resmodels.LoginPost{}
		userRows := db.ReadRows(dbmodels.User{}, rows)

		user := userRows[0].(dbmodels.User)
		returnRows = append(returnRows, resmodels.LoginPost{ID: user.ID, Username: user.Username, Auth: CreateUserToken(user.Username)})*/

	returnedRow = getUserLogin(w, reqmodels.LoginPost{Username: body.Username, Password: body.NewPassword})
	if returnedRow == nil {
		return
	}

	fmt.Println(returnedRow.Auth.Token)
	jio.WriteMsg(w, returnedRow, "")
}

// Logout handles the logout api call
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jio.Error(w, errors.New(methodNotAllowedMessage), 405)
		return
	}

	// Read body
	b := jio.Readbody(w, r)
	if b == nil {
		return
	}

	// Unmarshal body to body variable
	var body reqmodels.TokenPost
	err := json.Unmarshal(b, &body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
		errMsg := jio.CreateErrorMsgError(jio.InvalidJSON, []string{}, err)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	errMsg := app.DeleteUserToken(body.Username, body.AuthToken)
	if errMsg != nil {
		jio.WriteMsg(w, "", errMsg)
		return
	}

	jio.WriteMsg(w, "Succesfully logged out.", "")
}

// GetTokenInfo handles the token info api call
func GetTokenInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jio.Error(w, errors.New(methodNotAllowedMessage), 405)
		return
	}

	// Read body
	b := jio.Readbody(w, r)
	if b == nil {
		return
	}

	// Unmarshal body to body variable
	var body reqmodels.TokenPost
	err := json.Unmarshal(b, &body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
		errMsg := jio.CreateErrorMsgError(jio.InvalidJSON, []string{}, err)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	userToken, errMsg := app.GetUserToken(body.Username, body.AuthToken)

	if errMsg != nil {
		jio.WriteMsg(w, "", errMsg)
		return
	}

	jio.WriteMsg(w, userToken.GetTokenPostResponse(), "")
}

// RefreshTokenPost handles the refresh token api Call
func RefreshTokenPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jio.Error(w, errors.New(methodNotAllowedMessage), 405)
		return
	}

	// Read body
	b := jio.Readbody(w, r)
	if b == nil {
		return
	}

	// Unmarshal body to body variable
	var body reqmodels.TokenPost
	err := json.Unmarshal(b, &body)
	if err != nil {
		errMsg := jio.CreateErrorMsgError(jio.InvalidJSON, []string{}, err)
		jio.WriteMsg(w, "", errMsg)
		return
	}

	newToken, errMsg := app.RefreshUserToken(body.Username, body.AuthToken)
	if errMsg != nil {
		jio.WriteMsg(w, "", errMsg)
		return
	}

	userToken, errMsg := app.GetUserToken(body.Username, newToken)
	if errMsg != nil {
		jio.WriteMsg(w, "", errMsg)
		return
	}

	jio.WriteMsg(w, userToken.GetTokenPostResponse(), "")
}
