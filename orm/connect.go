package orm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"time"
)

var conn *gorm.DB

func Init() error {
	DbName := os.Getenv("DB_NAME")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", DbHost, DbUser, DbPassword, DbName, DbPort)

	var err error
	conn, err = gorm.Open(postgres.New(postgres.Config{
		DSN: DSN,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "Juu17Bot" + "_",
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return err
	}
	sqlDB, err := conn.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetConn() *gorm.DB {
	return conn
}
