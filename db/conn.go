package connections

import (
	configs "halo_suster/cfg"
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx as a stdlib driver for database/sql
	"github.com/jmoiron/sqlx"
)

func NewPgConn(config configs.Config) (*sqlx.DB, error) {
	ctx := context.Background()

	// Construct the connection string using the provided configuration
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DbUsername,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)

	// Connect to the database using sqlx
	db, err := sqlx.ConnectContext(ctx, "pgx", uri)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Set the connection pool configuration
	db.SetMaxOpenConns(90) // Adjust based on your database and workload
	db.SetMaxIdleConns(50) // Adjust based on your database and workload
	// db.SetConnMaxLifetime(1 * time.Hour)    // Set maximum connection lifetime
	// db.SetConnMaxIdleTime(30 * time.Minute) // Set maximum idle time for a connection

	// Test the connection before returning
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return db, nil
}
