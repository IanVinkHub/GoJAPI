package db

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"reflect"
	"strconv"
)

// QueryBuilder the builder for queries
type QueryBuilder struct {
	fromStr   []string
	selectStr []string
	whereStr  []string

	DBC Connection
}

// buildFrom builds the from of the query
func (qb *QueryBuilder) buildFrom() string {
	fromStr := "FROM "
	if len(qb.fromStr) == 0 {
		return ""
	}
	for i, v := range qb.fromStr {
		if i == 0 {
			fromStr += v
		} else {
			fromStr += ", " + v
		}
	}
	return fromStr
}

// buildSelect builds the select of the query
func (qb *QueryBuilder) buildSelect() string {
	selectStr := "SELECT "
	if len(qb.selectStr) == 0 {
		return "SELECT *"
	}
	for i, v := range qb.selectStr {
		if i == 0 {
			selectStr += v
		} else {
			selectStr += ", " + v
		}
	}
	return selectStr
}

// buildWhere builds the where of the query
func (qb *QueryBuilder) buildWhere() string {
	whereStr := "WHERE "
	if len(qb.whereStr) == 0 {
		return ""
	}
	for i, v := range qb.whereStr {
		if i == 0 {
			whereStr += v
		} else {
			whereStr += " AND " + v
		}
	}
	return whereStr
}

// ResetQuery resets the query
func (qb *QueryBuilder) ResetQuery() {
	qb.selectStr = []string{}
	qb.fromStr = []string{}
	qb.whereStr = []string{}
}

// All retrieves all data from selected table
func (qb *QueryBuilder) All() (*sql.Rows, error) {
	selectStr := qb.buildSelect()
	fromStr := qb.buildFrom()
	if fromStr == "" {
		return nil, errors.New("Table not selected")
	}
	whereStr := qb.buildWhere()

	qb.ResetQuery()

	err := qb.DBC.Open()
	if err != nil {
		return nil, err
	}
	defer qb.DBC.Close()

	rows, err := qb.DBC.conn.Query(selectStr + " " + fromStr + " " + whereStr)
	fmt.Println(selectStr + " " + fromStr + " " + whereStr)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Select adds to the select part of the query
func (qb *QueryBuilder) Select() *QueryBuilder {
	return qb
}

// From adds to the from part of the query
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.fromStr = append(qb.fromStr, table)
	return qb
}

// WhereRaw adds to the where part of the query
func (qb *QueryBuilder) WhereRaw(raw string) *QueryBuilder {
	qb.whereStr = append(qb.whereStr, raw)
	return qb
}

// WhereValue adds to the where part of the query by using column and value
func (qb *QueryBuilder) WhereValue(column string, value string) *QueryBuilder {
	value = "'" + html.EscapeString(value) + "'"
	qb.whereStr = append(qb.whereStr, column+" = "+value)
	return qb
}

// WhereValues adds to the where part of the query by using a interface
func (qb *QueryBuilder) WhereValues(values interface{}) *QueryBuilder {
	for i := 0; i < reflect.ValueOf(values).NumField(); i++ {
		val := reflectValToString(reflect.ValueOf(values).Field(i))
		qb.WhereValue(reflect.ValueOf(values).Type().Field(i).Name, val)
	}

	return qb
}

// GetWithID returns rows where id = given id from table
func (qb *QueryBuilder) GetWithID(id int64, table string) (*sql.Rows, error) {
	query := "SELECT * FROM " + table + " WHERE id=" + strconv.FormatInt(id, 10)
	return qb.DBC.conn.Query(query)
}

// Insert inserts values into table in database
func (qb *QueryBuilder) Insert(table string, values interface{}) (*sql.Rows, error) {
	query := "INSERT INTO " + table
	insertCols := "("
	insertVals := "("
	for i := 0; i < reflect.ValueOf(values).NumField(); i++ {
		val := html.EscapeString(reflectValToString(reflect.ValueOf(values).Field(i)))
		if i > 0 {
			insertCols += ", "
			insertVals += ", "
		}
		insertCols += reflect.ValueOf(values).Type().Field(i).Name
		insertVals += "'" + val + "'"
	}
	insertCols += ")"
	insertVals += ")"

	query += insertCols + " VALUES " + insertVals

	err := qb.DBC.Open()
	if err != nil {
		return nil, err
	}
	defer qb.DBC.Close()

	result, err := qb.DBC.conn.Exec(query)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return qb.GetWithID(id, table)
}

// Update updates values from table in database
func (qb *QueryBuilder) Update(table string, values interface{}) (*sql.Rows, error) {
	query := "UPDATE " + table + " SET "

	for i := 0; i < reflect.ValueOf(values).NumField(); i++ {
		val := html.EscapeString(reflectValToString(reflect.ValueOf(values).Field(i)))
		if i > 0 {
			query += ", "
		}
		query += reflect.ValueOf(values).Type().Field(i).Name
		query += " = "
		query += "'" + val + "'"
	}
	query += " " + qb.buildWhere()

	err := qb.DBC.Open()
	if err != nil {
		return nil, err
	}
	defer qb.DBC.Close()

	fmt.Println(query)
	_, err = qb.DBC.conn.Exec(query)
	if err != nil {
		return nil, err
	}

	qb.ResetQuery()

	return qb.From(table).WhereValues(values).All()
}

func reflectValToString(value reflect.Value) string {
	switch value.Type() {
	case reflect.TypeOf(0):
		return strconv.Itoa(value.Interface().(int))
	}
	return value.String()
}
