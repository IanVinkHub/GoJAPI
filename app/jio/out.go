package jio // JSON IN OUT

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Const for custom errors
const (
	InvalidParameters   int = 4001
	InvalidJSON         int = 4002
	CredentialsExist    int = 6001
	CredentialsNotFound int = 6002
	TokenExpired        int = 6003
	UserTokenNotFound   int = 6004
)

// ErrorMsg the message to send to the client if an error occures
type ErrorMsg struct {
	Code       int      `json:"code"`
	CustomCode int      `json:"custom_code"`
	ErrorType  string   `json:"type"`
	ErrorParam []string `json:"param"`
	Message    string   `json:"message"`
}

func getErrorMessage(code int) string {
	switch code {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 500:
		return "Internal server error"
	}

	return ""
}

// CreateErrorMsgError creates error message with given cusCode, params and error
func CreateErrorMsgError(cusCode int, errorParam []string, err error) ErrorMsg {
	var errorType string
	message := err.Error()
	code := 500

	if errorParam == nil {
		errorParam = make([]string, 0)
	}

	switch cusCode {
	case InvalidParameters:
		code = 400
		errorType = "Invalid parameters"
		message = "The following parameters does not exist or empty: ["
		for i, elm := range errorParam {
			message += elm
			if i < len(errorParam)-1 {
				message += ", "
			}
		}
		message += "]"
		break
	case InvalidJSON:
		code = 400
		errorType = "Invalid JSON"
		break
	case CredentialsExist:
		code = 401
		message = "Username already in use"
		break
	case CredentialsNotFound:
		code = 401
		message = "No user found with these credentials"
		break
	case TokenExpired:
		code = 401
		message = "Token is expired please refresh the token using the /auth/refresh endpoint."
	case UserTokenNotFound:
		code = 401
		message = "Token or user could not be found please login."
	}

	if errorType == "" {
		errorType = getErrorMessage(code)
	}

	return ErrorMsg{code, cusCode, errorType, errorParam, message}
}

// CreateErrorMsg create error message with given cusCode and params
func CreateErrorMsg(cusCode int, errorParam []string) ErrorMsg {
	return CreateErrorMsgError(cusCode, errorParam, errors.New(""))
}

// Error outputs error
func Error(w http.ResponseWriter, newErr error, code int) {
	errorMsg := ErrorMsg{code, -1, "", make([]string, 0), newErr.Error()}

	errorMsg.ErrorType = getErrorMessage(code)

	output, err := CreateMsg("error", "", errorMsg)
	if err != nil {
		Error(w, err, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(output)
}

// CreateMsg creates an unified JSON structure for writing
func CreateMsg(status string, dataMsg interface{}, errMsg interface{}) ([]byte, error) {
	outputStruct := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
		Error  interface{} `json:"error"`
	}{status, dataMsg, errMsg}

	output, err := json.Marshal(outputStruct)
	return output, err
}

// WriteMsg writes a unified JSON structure to the client
func WriteMsg(w http.ResponseWriter, dataMsg interface{}, errMsg interface{}) {
	status := "succes"
	if errMsg != "" {
		status = "error"
	}

	output, err := CreateMsg(status, dataMsg, errMsg)
	if err != nil {
		Error(w, err, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
