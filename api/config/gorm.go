package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(viper *viper.Viper) *gorm.DB {
	username := os.Getenv("username")
	password := os.Getenv("password")
	host := os.Getenv("host")
	port := os.Getenv("port")
	database := os.Getenv("name")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, username, password, database, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
		// 	SlowThreshold:             time.Second * 5,
		// 	Colorful:                  false,
		// 	IgnoreRecordNotFoundError: true,
		// 	ParameterizedQueries:      true,
		// 	LogLevel:                  logger.Info,
		// }),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	_, err = db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// connection.SetMaxIdleConns(idleConnection)
	// connection.SetMaxOpenConns(maxConnection)
	// connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}
