package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb(user string, pass string, hostname string, port string) {
	var err error
	dsn := user + ":" + pass + "@tcp(" + hostname + ":" + port + ")/policy?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Failed to connect to Db")
	}
}
