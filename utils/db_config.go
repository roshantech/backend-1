package utils

import (
	model "backend/models"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	dbMux sync.Mutex
	once  sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		initDB()
	})
	return db
}

func CloseDB() {
	dbMux.Lock()
	defer dbMux.Unlock()
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("Error getting underlying database:", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Println("Error closing the database connection:", err)
		}
	}
}

func initDB() {
	var err error
	db, err = gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/test"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = model.MigrateModels(db)
	if err != nil {
		log.Fatal("Failed to perform auto migration:", err)
	}
}
