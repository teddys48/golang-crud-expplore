package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidateRequest(r *http.Request, validate *validator.Validate, v any) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}

	err = validate.Struct(v)
	if err != nil {
		var errMessage string
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
			fieldName := err.Field()
			field, _ := reflect.TypeOf(v).Elem().FieldByName(fieldName)
			jsonField, _ := field.Tag.Lookup("json")
			errMessage = jsonField + " is " + err.ActualTag()
		}

		return fmt.Errorf("%v", errMessage)
	}

	return nil
}
