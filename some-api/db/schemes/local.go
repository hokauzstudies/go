package db

import db "pep-api/db/tools"

var localCollection = "local"
var localLocalRelationCollection = "local_has_local"

// AddLocal -
func AddLocal(u interface{}) (interface{}, error) {

	res, err := db.ExecuteQuery("insert", localCollection, u, nil, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ExistsLocal -
func ExistsLocal(where map[string]interface{}) (bool, error) {
	res, err := db.ExecuteQuery("exists", localCollection, map[string]interface{}{"id": true}, where, nil)
	if err != nil {
		return true, err
	}
	return res.(bool), nil
}

// LocalHasLocal -
func LocalHasLocal(where map[string]interface{}) (bool, error) {
	res, err := db.ExecuteQuery("exists", localLocalRelationCollection, map[string]interface{}{"id": true}, where, nil)
	if err != nil {
		return true, err
	}
	return res.(bool), nil
}

// GetLocalByID -
func GetLocalByID(id int) (interface{}, error) {
	res, err := db.ExecuteQuery("select", localCollection, map[string]interface{}{"*": true}, map[string]interface{}{"id": id}, nil)
	if err != nil {
		return nil, err
	}
	return res.([]interface{})[0], nil
}

// GetLocals -
func GetLocals(argsIn ...interface{}) ([]interface{}, error) {
	// paginator := argsIn[1]

	res, err := db.ExecuteQuery("select", localCollection, map[string]interface{}{"*": true}, argsIn[0], nil)
	if err != nil {
		return nil, err
	}
	return res.([]interface{}), nil
}

// UpdateLocal -
func UpdateLocal(id int, u interface{}) (interface{}, error) {

	res, err := db.ExecuteQuery("update", localCollection, u, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	return res, nil
}
