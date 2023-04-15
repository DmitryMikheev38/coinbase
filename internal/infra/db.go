package infra

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConfigDB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func ConnectToDB(db ConfigDB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.User, db.Password, db.Host, db.Port, db.Name)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
