package models

import (
	"fmt"
	"santapKlik/configs"

	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitModel(c configs.ProgramConfig) *gorm.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect database, ", err.Error())
	}

	return db
}
// func InitModel(config configs.ProgramConfig) *gorm.DB {
// 	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		logrus.Error("Model : cannot connect to database, ", err.Error())
// 		return nil
// 	}

// 	return db
// }

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&Jajanan{})
	db.AutoMigrate(&Makanan{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Order{})
}
