package postgres

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection(host string, port int, dbName, user, password, timezone string) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		host,
		user,
		password,
		dbName,
		port,
		timezone,
	)

	logrus.Info("Connecting to Postgres: ", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logrus.Info("Connected to Postgres: ", dsn)
	return db
}
