package jio // JSON IN OUT
import (
	"io/ioutil"
	"net/http"
	"reflect"
)

// Readbody reads body to byte array and if failed return error 500
func Readbody(w http.ResponseWriter, r *http.Request) []byte {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		Error(w, err, 500)
		return nil
	}

	return b
}

// CheckParamsStruc checks if all parameters in the structure are set
func CheckParamsStruc(body interface{}) []string {
	var errorParam []string
	reflectedBody := reflect.ValueOf(body)

	for i := 0; i < reflectedBody.NumField(); i++ {
		field := reflectedBody.Type().Field(i)
		reflectVal := reflectedBody.FieldByName(field.Name)
		if !reflectVal.IsValid() || reflectVal.String() == "" {
			errorParam = append(errorParam, field.Tag.Get("json"))
		}
	}

	return errorParam
}

// CheckParams checks if all listed params in params are set in body struc
func CheckParams(body interface{}, params []string) []string {
	var errorParam []string
	reflectedBody := reflect.ValueOf(body)

	for _, val := range params {
		reflectVal := reflectedBody.FieldByName(val)
		if !reflectVal.IsValid() || reflectVal.String() == "" {
			errorParam = append(errorParam, val)
		}
	}

	return errorParam
}
