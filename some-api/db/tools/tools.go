package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"pep-api/db"
	"reflect"
	"strconv"
	"strings"
)

var relationalTables = map[string]interface{}{
	"local_id": "user_local",
}

var statementModel = map[string]string{
	"update": "UPDATE %s SET %s",
	"insert": "INSERT INTO %s (%s) \nVALUES\n (%s)",
	"select": "SELECT %s FROM %s",
	"exists": "SELECT %s FROM %s"}

//ExecuteQuery -
//	Tipo de Query 'update' | 'insert' | 'select' | 'exists'
//	Nome da Tabela string
//	Dados []interface{}
//	WHERE map[string]interface{} || nil
//	Dados especiais []string
func ExecuteQuery(argsIn ...interface{}) (res interface{}, err error) {

	argsIn[2] = structToMap(argsIn[2])

	sqlStatement := mountStatement(argsIn...)
	fmt.Println(sqlStatement[0])
	fmt.Println(sqlStatement[1])
	fmt.Println(sqlStatement[2])
	var stmt *sql.Stmt
	var rows *sql.Rows
	var result sql.Result
	stmt, err = db.Connection.Prepare(sqlStatement[0].(string))
	defer stmt.Close()
	switch argsIn[0].(string) {
	case "exists":
		rows, err := stmt.Query(sqlStatement[1].([]interface{})...)
		if err != nil {
			return nil, err
		}
		res, err = existsQuery(rows)
	case "select":
		fmt.Println(sqlStatement[0])
		fmt.Println(sqlStatement[1])
		rows, err = stmt.Query(sqlStatement[1].([]interface{})...)
		if err != nil {
			return nil, err
		}
		res, err = selectQuery(rows)
	case "update":
		_, err = stmt.Exec(sqlStatement[1].([]interface{})...)
		if err != nil {
			return nil, err
		}
		res = argsIn[2]
	case "insert":

		result, err = stmt.Exec(sqlStatement[1].([]interface{})...)
		if err != nil {
			return nil, err
		}
		insertedID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		intID := strconv.Itoa(int(insertedID))

		////Insertion to relational tables
		if sqlStatement[2] != nil {

			for specialIndex, specialArr := range sqlStatement[2].(map[string]interface{}) {

				switch specialArr.(type) {
				case []interface{}:
					if len(specialArr.([]interface{})) > 0 {
						for _, specialData := range specialArr.([]interface{}) {

							specialSQLStatement := mountStatement("insert", relationalTables[specialIndex], map[string]interface{}{"user_id": intID, specialIndex: specialData}, nil)
							_, err := db.Connection.Exec(specialSQLStatement[0].(string), specialSQLStatement[1].([]interface{})...)
							if err != nil {
								log.Println(err.Error())
							}
						}
					}
				case float64, string, int:

					specialSQLStatement := mountStatement("insert", relationalTables[specialIndex], map[string]interface{}{"user_id": intID, specialIndex: specialArr}, nil)
					_, err := db.Connection.Exec(specialSQLStatement[0].(string), specialSQLStatement[1].([]interface{})...)
					if err != nil {
						log.Println(err.Error())
					}
				default:
					fmt.Println(reflect.TypeOf(specialArr))
				}

				// fmt.Println(newSpecialData)
			}
		}

		argsIn[2].(map[string]interface{})["id"] = intID

		res = argsIn[2]
	}

	if err != nil {
		return nil, err
	}
	return res, nil
}

func existsQuery(rows *sql.Rows) (res bool, err error) {
	var count int
	if err != nil {
		return
	}
	count, err = checkCount(rows)
	if err != nil {
		return
	}

	res = false
	if count > 0 {
		res = true
	}
	return
}

func selectQuery(rows *sql.Rows) (res []interface{}, err error) {

	if rows != nil {
		columns, _ := rows.Columns()
		count := len(columns)

		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		for rows.Next() {
			row := make(map[string]interface{}, count)
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			rows.Scan(valuePtrs...)

			for i, col := range columns {
				val := values[i]
				b, ok := val.([]byte)
				var v interface{}
				if ok {
					v = string(b)
				} else {
					v = val
				}
				row[col] = v
			}

			res = append(res, row)
		}
		return
	}
	return
}

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return
		}
	}
	return
}

func mountStatement(argsIn ...interface{}) (argsOut []interface{}) {

	var rawWhereValues []interface{}
	var specialData = make(map[string]interface{})
	var whereStatement string

	statementType := argsIn[0]
	tableName := argsIn[1]
	data := argsIn[2].(map[string]interface{})
	queryString := statementModel[statementType.(string)]

	if argsIn[3] != nil && len(argsIn) >= 4 {
		var whereStatementArr []string
		//Where just for equality
		for column, value := range argsIn[3].(map[string]interface{}) {
			rawWhereValues = append(rawWhereValues, value)
			whereStatementArr = append(whereStatementArr, fmt.Sprintf("%s = ?", column))
		}
		whereStatement = strings.Join(whereStatementArr, ", ")
		queryString += " WHERE %s"
	}

	if len(argsIn) == 5 && argsIn[4] != nil {
		fmt.Println("Argsin 4: ", argsIn[4])
		for _, specialKey := range argsIn[4].([]string) {
			fmt.Println("Special Key: ", specialKey)
			specialData[specialKey] = data[specialKey]
			delete(data, specialKey)
		}
	}
	fmt.Println("Data: ", data)

	queryArguments, rawValues := getQueryArguments(data, statementType.(string), tableName.(string))
	if argsIn[3] != nil && len(argsIn) >= 4 {
		queryArguments = append(queryArguments, whereStatement)
		rawValues = append(rawValues, rawWhereValues...)
	}

	sqlStatement := fmt.Sprintf(queryString, queryArguments...)
	argsOut = append(argsOut, sqlStatement)
	argsOut = append(argsOut, rawValues)
	argsOut = append(argsOut, specialData)

	return
}

func getQueryArguments(keys map[string]interface{}, statementType, tableName string) (queryArguments []interface{}, rawValues []interface{}) {
	valueArr, columnArr := lookupInfo(keys)
	var bindArr []string

	for _, val := range valueArr {
		rawValues = append(rawValues, val)
		bindArr = append(bindArr, "?")
	}

	switch statementType {
	case "update":
		var updateStatementArr []string
		for _, column := range columnArr {
			updateStatementArr = append(updateStatementArr, fmt.Sprintf("%s = ?", column))
		}
		updateStatement := strings.Join(updateStatementArr, ", ")
		queryArguments = append(queryArguments, tableName, updateStatement)
	case "insert":
		values := strings.Join(bindArr[:], ", ")
		columns := strings.Join(columnArr[:], ", ")
		queryArguments = append(queryArguments, tableName, columns, values)
	case "select", "exists":
		columns := strings.Join(columnArr[:], ", ")
		queryArguments = append(queryArguments, columns, tableName)
	}
	return
}

func lookupInfo(keys map[string]interface{}) (valueString, rawColumns []string) {

	for key, element := range keys {
		rawColumns = append(rawColumns, key)
		switch element.(type) {
		case string:
			valueString = append(valueString, fmt.Sprintf("%s", element))
		case int:
			valueString = append(valueString, fmt.Sprintf("%d", element))
		case float64:
			valueString = append(valueString, fmt.Sprintf("%0.2f", element))
		}
	}
	return
}

func structToMap(data interface{}) (keys map[string]interface{}) {
	inrec, _ := json.Marshal(data)
	json.Unmarshal(inrec, &keys)
	return
}
