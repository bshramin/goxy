package postgres

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConnection(host string, port int, dbName, user, password string) (*gorm.DB, error) {
	dsn := url.URL{
		User:     url.UserPassword(user, password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", host, port),
		Path:     dbName,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})

	return db, err
}

func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Ping()
	return err
}
