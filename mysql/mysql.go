package mysql

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConnection(host string, port int, dbName, user, password string) *gorm.DB {
	dsn := url.URL{
		User:     url.UserPassword(user, password),
		Scheme:   "mysql",
		Host:     fmt.Sprintf("%s:%d", host, port),
		Path:     dbName,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	logrus.Info("Connecting to MySQL: ", dsn.String())

	db, err := gorm.Open(mysql.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logrus.Info("Connected to MySQL: ", dsn.String())
	return db
}

func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Ping()
	return err
}
