package main

import (
	"fmt"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/health", service.HealthHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/api/login", service.LoginHandler)
	http.HandleFunc("/api/user", service.ProfileHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
