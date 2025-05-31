package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Westlam/S13/config"
	"github.com/Westlam/S13/pkg/api"
)

func Run(config config.Config, db *sql.DB) {
	fs := http.FileServer(http.Dir("./web"))
	//fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", http.StripPrefix("/", fs))
	api.Init(db)

	port := config.TODO_PORT
	log.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
