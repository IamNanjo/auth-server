package main

import (
	"embed"
	"github.com/IamNanjo/authserver/backend"
	"github.com/IamNanjo/authserver/db"
	"os"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	addr := os.Getenv("AUTH_SERVER_PORT")
	if addr == "" {
		addr = ":8080"
	}

	db.Initialize()

	os.Stdout.WriteString("Listening on " + addr + "\n")
	backend.StartServer(addr, staticFiles)
}
