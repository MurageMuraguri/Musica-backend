package db

import (
	"musica/utils"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func Db() *gorm.DB {

	db, err = gorm.Open("postgres", utils.Env("POSTGRES"))
	if err != nil {
		log.Println(err)
	}
	return db
}
