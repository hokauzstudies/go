package db

import db "pep-api/db/tools"

var userCollection = "user"
var userLocalRelationCollection = "user_has_local"

// AddUser -
func AddUser(u interface{}) (interface{}, error) {

	res, err := db.ExecuteQuery("insert", userCollection, u, nil, []string{"locals", "local_id"})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ExistsUser -
func ExistsUser(where map[string]interface{}) (bool, error) {
	res, err := db.ExecuteQuery("exists", userCollection, map[string]interface{}{"id": true}, where, nil)
	if err != nil {
		return true, err
	}
	return res.(bool), nil
}

// UserHasLocal -
func UserHasLocal(where map[string]interface{}) (bool, error) {
	res, err := db.ExecuteQuery("exists", userLocalRelationCollection, map[string]interface{}{"id": true}, where, nil)
	if err != nil {
		return true, err
	}
	return res.(bool), nil
}

// GetUserHash -
func GetUserHash(email string, pass string) (interface{}, error) {
	res, err := db.ExecuteQuery("select", userCollection, map[string]interface{}{"password": true}, map[string]interface{}{"email": email}, nil)
	if err != nil {
		return nil, err
	}
	if len(res.([]interface{})) > 0 {
		return res.([]interface{})[0], nil
	}
	return nil, nil
}

// GetUserByID -
func GetUserByID(id int) (interface{}, error) {
	res, err := db.ExecuteQuery("select", userCollection, map[string]interface{}{"*": true}, map[string]interface{}{"id": id}, nil)
	if err != nil {
		return nil, err
	}
	if len(res.([]interface{})) > 0 {
		return res.([]interface{})[0], nil
	}
	return nil, nil
}

// GetUsers -
func GetUsers(argsIn ...interface{}) ([]interface{}, error) {
	// paginator := argsIn[1]

	res, err := db.ExecuteQuery("select", userCollection, map[string]interface{}{"*": true}, argsIn[0], nil)
	if err != nil {
		return nil, err
	}
	return res.([]interface{}), nil
}

// GetUserByMPI - // TODO be removed
func GetUserByMPI(mpi string, cnes string) (interface{}, error) {
	return nil, nil
}

// UpdateUser -
func UpdateUser(id int, u interface{}) (interface{}, error) {

	res, err := db.ExecuteQuery("update", userCollection, u, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CheckPass -
func CheckPass(id int, pass string) (bool, error) {
	return true, nil
}
