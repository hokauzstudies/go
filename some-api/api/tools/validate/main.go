package validate

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func NewValidator(data interface{}) (errStr string) {
	validate = validator.New()
	err := validate.Struct(data)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return "Nil validators"
		}

		var arrErr []string
		for _, err := range err.(validator.ValidationErrors) {
			arrErr = append(arrErr, "Missing field "+err.Field()+" with type "+err.Kind().String())
			// fmt.Println(err.Namespace()) // can differ when a custom TagNameFunc is registered or
			// fmt.Println(err.Field())     // by passing alt name to ReportError like below
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
		}
		errStr = strings.Join(arrErr, ", ")
		// from here you can create your own error messages in whatever language you wish
		return
	}
	return
}
