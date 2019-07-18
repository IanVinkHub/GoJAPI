package controllers

import (
	"fmt"

	"github.com/IanVinkHub/gojapi/app/db"
)

var qb db.QueryBuilder

func init() {
	qb.DBC = db.Connection{Database: "gojapi"}
	qb.DBC.Default()
	err := qb.DBC.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
