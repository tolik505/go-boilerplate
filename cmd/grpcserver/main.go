package main

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"goboilerplate/pkg/storage/db"
	"net"
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
	grpcApp, err := InitializeGRPCApp(dbConf)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Couldn't initialise GRPC app"))
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		err = errors.Wrap(err, "failed to listen tcp")
		log.Fatal(err)
	}

	go grpcApp.Run(lis)

	log.Info("gRPC app is running")

	grpcApp.WaitForInterrupt()

	// Gracefully stop important processes if any

	log.Info("exited")
}
