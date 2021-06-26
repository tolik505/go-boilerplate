package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"goboilerplate/pkg/api/minio"
	"goboilerplate/pkg/httpapp"
	"goboilerplate/pkg/storage/db"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	dbConf := db.Config{
		Host:   os.Getenv("DB_HOST"),
		Port:   os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
		User:   os.Getenv("DB_USER"),
		Pass:   os.Getenv("DB_PASS"),
	}
	minioConf := minio.Config{
		Host:   os.Getenv("MINIO_HOST"),
		Key:    os.Getenv("MINIO_KEY"),
		Secret: os.Getenv("MINIO_SECRET"),
	}
	bucket := minio.Bucket(os.Getenv("MINIO_BUCKET"))
	httpPort := httpapp.Port(os.Getenv("HTTP_PORT"))
	httpApp, err := InitializeHTTPApp(httpPort, dbConf, minioConf, bucket)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Couldn't initialise HTTP app"))
	}

	go httpApp.Run()

	httpApp.WaitForInterrupt()

	// Gracefully stop important processes if any

	log.Info("exited")
}
