package migrate

import (
	"musica/model"
	"musica/utils"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error
type UserIndex struct {
	Name  string `gorm:"check:name_checker,name <> 'jinzhu'"`

  }
func Migrate() {
	db, err = gorm.Open("postgres", utils.Env("POSTGRES"))
	//db.AutoMigrate(&model.Driver{},&model.Licence{})
	db.Debug().AutoMigrate(
		&model.Audio{},

		&model.User{},

		&model.File{},
		&model.Login{},

		&UserIndex{},
	)
	if err != nil {
		log.Println(err)
	}
}
