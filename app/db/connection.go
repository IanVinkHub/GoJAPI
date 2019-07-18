package db

import (
	"database/sql"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql" // Makes sure it is using this mysql driver
)

// Connection is the connection to the database
type Connection struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string

	conn *sql.DB
}

// Default sets the default options for the Connection
func (dbc *Connection) Default() {
	if dbc.Host == "" {
		dbc.Host = "localhost"
	}

	if dbc.Port == "" {
		dbc.Port = "3306"
	}

	if dbc.Username == "" {
		dbc.Username = "root"
	}
}

// Open opens a connection to the database
func (dbc *Connection) Open() error {
	pw := ""
	if dbc.Password != "" {
		pw = ":" + dbc.Password
	}
	//dbc.Username+pw
	conn, err := sql.Open("mysql", dbc.Username+pw+"@tcp("+dbc.Host+":"+dbc.Port+")/"+dbc.Database)
	if err != nil {
		return err
	}

	if err = conn.Ping(); err != nil {
		return err
	}

	dbc.conn = conn
	return nil
}

// Close closes a connection to the database
func (dbc *Connection) Close() error {
	err := dbc.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

// ReadRows into array of interface type
func ReadRows(struc interface{}, rows *sql.Rows) []interface{} {
	var structs []interface{}
	// Read all rows into structs
	for rows.Next() {
		columns := make([]interface{}, 0)
		for i := 0; i < reflect.TypeOf(struc).NumField(); i++ {
			// Switch trough all reflect types and set type (interface{} is default)
			switch reflect.TypeOf(struc).Field(i).Type {
			case reflect.TypeOf(""):
				field := reflect.New(reflect.TypeOf(struc).Field(i).Type).Elem().Interface().(string)
				columns = append(columns, &field)
			case reflect.TypeOf(0):
				field := reflect.New(reflect.TypeOf(struc).Field(i).Type).Elem().Interface().(int)
				columns = append(columns, &field)
			default:
				field := reflect.New(reflect.TypeOf(struc).Field(i).Type).Elem().Interface()
				columns = append(columns, &field)
			}
		}

		if err := rows.Scan(columns...); err != nil {
			log.Fatal(err)
		}
		structs = append(structs, columns)
	}

	// Create and format strucs for returning
	var retStructs []interface{}
	for _, v := range structs {
		// Create new struc of struc type then convert interface{} to []interface{}
		newStruc := reflect.New(reflect.TypeOf(struc)).Elem()
		strucArr := reflect.ValueOf(v).Interface().([]interface{})

		// Loop trough all fields in struc
		for j := 0; j < reflect.TypeOf(struc).NumField(); j++ {
			// select current fieldval and set struc to val
			fieldVal := reflect.ValueOf(strucArr[j]).Elem()
			newStruc.FieldByName(reflect.TypeOf(struc).Field(j).Name).Set(fieldVal)
		}
		retStructs = append(retStructs, newStruc.Interface())
	}
	return retStructs
}
