package app

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/IanVinkHub/gojapi/app/jio"
	"github.com/IanVinkHub/gojapi/models/resmodels"
)

// Default time format for javascript is RFC2822 wich is "Mon, 02 Jan 2006 15:04:05 +0200"
// We use this for easy conversion in the browser
const timeFormat = "Mon, 02 Jan 2006 15:04:05 +0200"

// User is a struct for every logged in client
type User struct {
	AuthToken         string
	Username          string
	TokenCreated      time.Time
	TokenExpiresAt    time.Time
	TokenGetDeletedAt time.Time
	Timer             *time.Timer
}

const letters = "abdefghijklmnqrstuvxyzABDEFGHIJKLMNQRSTUVXYZ0123456789!$&'()*+-.:@_~"

var (
	// Users contains all the sessions with sessionkey as key
	Users map[string]User
	// CreateSessionMux is the mux to keep everything safe
	CreateSessionMux sync.Mutex
)

func init() {
	// Create the User map
	Users = make(map[string]User)
}

// CreateUserToken creates a user bound to a token
func CreateUserToken(username string) string {
	var token string
	generated := false
	tokenbytes := make([]byte, 64)
	for !generated {
		_, err := rand.Read(tokenbytes)
		if err != nil {
			// handle error here
		}

		for i, b := range tokenbytes {
			tokenbytes[i] = letters[b%byte(len(letters))]
		}
		token = string(tokenbytes)

		if _, ok := Users[token]; !ok {
			now := time.Now()

			timer := time.AfterFunc(time.Second*60, func() {
				fmt.Println("testing some stuff with token " + token)
				delete(Users, token)
			})

			Users[token] = User{
				AuthToken:         token,
				Username:          username,
				TokenCreated:      now,
				TokenExpiresAt:    now.Add(time.Second * 40),
				TokenGetDeletedAt: now.Add(time.Second * 60),
				Timer:             timer,
			}

			generated = true
		}
	}

	return token
}

// GetUserToken gets the token bound to the user if exists
func GetUserToken(username string, token string) (*User, *jio.ErrorMsg) {
	if user, ok := Users[token]; ok && user.Username == username {

		if time.Now().After(user.TokenExpiresAt) {
			errMsg := jio.CreateErrorMsg(jio.TokenExpired, []string{})
			return &user, &errMsg
		}

		return &user, nil
	}

	errMsg := jio.CreateErrorMsg(jio.UserTokenNotFound, []string{})
	return nil, &errMsg
}

// RefreshUserToken refreshes the user token
func RefreshUserToken(username string, token string) (string, *jio.ErrorMsg) {
	errMsg := DeleteUserToken(username, token)
	if errMsg != nil && errMsg.CustomCode == jio.UserTokenNotFound {
		return "", errMsg
	}

	return CreateUserToken(username), nil
}

// DeleteUserToken deletes the user out of the Users map
func DeleteUserToken(username string, token string) *jio.ErrorMsg {
	user, errMsg := GetUserToken(username, token)
	if errMsg != nil && errMsg.CustomCode == jio.UserTokenNotFound {
		return errMsg
	}
	user.Timer.Stop()
	delete(Users, token)

	return nil
}

// GetTokenPostResponse returns the Post data for Get Token
func (user *User) GetTokenPostResponse() resmodels.TokenPost {
	return resmodels.TokenPost{Token: user.AuthToken, CreatedAt: user.TokenCreated.Format(timeFormat), ExpireAt: user.TokenExpiresAt.Format(timeFormat), DeleteAt: user.TokenGetDeletedAt.Format(timeFormat)}
}
