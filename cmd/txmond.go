package main

import (
	"log"
	"os"

	"github.com/AndreiRegiani/txmond/cmd/daemon"
	"github.com/AndreiRegiani/txmond/cmd/rest"
	"github.com/AndreiRegiani/txmond/cmd/storage"
	"github.com/AndreiRegiani/txmond/cmd/utils"
)

func main() {
	if err := utils.LoadEnvFile(".env"); err != nil {
		log.Fatal("txmond: error loading .env file")
	}

	storageBackend := os.Getenv("TXMOND_STORAGE_BACKEND")

	switch storageBackend {
	case "temp":
		log.Println("txmond: using 'temp' storage")
		storage.Db = storage.NewTempStorage()
	case "redis":
		log.Println("txmond: using 'redis' storage")
		storage.Db = storage.NewRedisStorage()
	default:
		log.Fatal("txmond: invalid storage backend")
	}

	// thread: REST API interface
	go rest.Listen()

	// thread: background worker
	daemon.Start()
}
