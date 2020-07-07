package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, _ := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
}
