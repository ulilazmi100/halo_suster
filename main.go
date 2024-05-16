package main

import (
	"log"

	configs "halo_suster/cfg"
	connections "halo_suster/db"
	"halo_suster/server"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Set up database connection pool
	dbPool, err := connections.NewPgConn(config)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbPool.Close()

	s := server.NewServer(dbPool)

	s.RegisterRoute(config)

	s.Run(config)
}
