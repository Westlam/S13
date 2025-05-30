package main

import (
	"log"

	"github.com/Westlam/S13/config"
	"github.com/Westlam/S13/pkg/db"
	"github.com/Westlam/S13/pkg/server"
)

func main() {
	config := config.New()

	db, err := db.Init(config)
	if err != nil {
		log.Fatal("Error creating database:", err)
	}
	defer db.Close()

	server.Run(config, db)
}
