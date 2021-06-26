package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host   string
	Port   string
	DBName string
	User   string
	Pass   string
}

// InitDB Opens new connection with Mysql
func InitDB(conf Config) (*gorm.DB, error) {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.DBName,
	)

	return gorm.Open(gormmysql.Open(dataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
