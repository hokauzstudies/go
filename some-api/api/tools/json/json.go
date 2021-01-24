package json

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
)

// Stringify - object to json bytes
func Stringify(data interface{}) ([]byte, error) {
	json, err := jsoniter.Marshal(data)

	if err != nil {
		return nil, err
	}
	if !jsoniter.Valid(json) {
		return nil, errors.New("invalid JSON data")
	}

	return json, nil
}

//Parse - json bytes to object
func Parse(data []byte, model interface{}) error {
	if !jsoniter.Valid(data) {
		return errors.New("invalid JSON data")
	}

	var err = jsoniter.Unmarshal(data, model)
	if err != nil {
		return err
	}
	return nil
}
