package initializers

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb(user string, pass string, hostname string, port string) {
	var err error
	dsn := user + ":" + pass + "@tcp(" + hostname + ":" + port + ")/goexcel?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		// NamingStrategy: schema.NamingStrategy{
		// 	SingularTable: true,
		// },
	})

	if err != nil {
		panic("Failed to connect to Db")
	}

	// Configure the connection pool
	sqlDB, err := DB.DB() // Get the underlying sql.DB object
	if err != nil {
		panic("Failed to get database instance")
	}

	// Set connection pool configurations
	sqlDB.SetMaxOpenConns(100)          // Set the maximum number of open connections
	sqlDB.SetMaxIdleConns(10)           // Set the maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set the maximum connection lifetime
}
