package testhelper

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestDB contains DB instance
var TestDB *gorm.DB

type DBName string

// InitTestDB Opens new connection with test Mysql
func InitTestDB() {
	if TestDB != nil {
		return
	}

	dbName := os.Getenv("TEST_DB_NAME")

	if dbName == "" {
		log.Fatal("TEST_DB_NAME env var is missing")
	}
	TestDB = InitSpecificTestDB(DBName(dbName))
}

// InitSpecificTestDB Opens a new connection for a specific db
func InitSpecificTestDB(dbName DBName) *gorm.DB {
	host := os.Getenv("TEST_DB_HOST")
	port := os.Getenv("TEST_DB_PORT")

	dataSource := fmt.Sprintf("root:password@tcp(%s:%s)/", host, port)
	db, err := gorm.Open(gormmysql.Open(dataSource), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	if err := db.Exec(fmt.Sprintf("CREATE database IF NOT EXISTS %s", dbName)).Error; err != nil {
		log.Fatal(err)
	}

	dataSource = fmt.Sprintf("root:password@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", host, port, dbName)
	db, err = gorm.Open(gormmysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	_, b, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join("file://", filepath.Dir(b), "../../migrations")
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		string(dbName),
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	return db
}

func BrokenDB() *gorm.DB {
	dbMock, _, _ := sqlmock.New()
	db, _ := gorm.Open(gormmysql.New(gormmysql.Config{
		Conn: dbMock,
	}), &gorm.Config{})

	return db
}

func SeedFixtures(t *testing.T, db *gorm.DB, path string) {
	sqlDB, _ := db.DB()
	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB), // You database connection
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(path), // the directory containing the YAML files
		testfixtures.Location(time.UTC),
	)
	if err != nil {
		t.Fatal("couldn't instantiate fixtures", err)
	}
	if err = fixtures.Load(); err != nil {
		t.Fatal("couldn't load fixtures", err)
	}
}

func CleanDB(db *gorm.DB) error {
	tables := []string{
		"posts",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s; ALTER TABLE %[1]s AUTO_INCREMENT = 1", table)).Error; err != nil {
			return errors.Errorf("Couldn't delete from table %s: %v", table, err)
		}
	}

	return nil
}
