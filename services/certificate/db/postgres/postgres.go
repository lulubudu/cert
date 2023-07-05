package postgres

import (
	"certificate/db"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	host     = "postgres"
	port     = 5432
	username = "docker"
	password = "docker"
	dbname   = "certificate_dev"
)

// Postgres composites sql.DB and represents a postgres connection.
type Postgres struct {
	*sql.DB
}

// Connect returns a Postgres with an active connection, and sets up graceful
// shutdown for the database connection.
func Connect() (db.Database, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	// create connection to postgres
	sqlDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// ping db to check the connection - postgres doesn't check liveness in `Open`
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	// handle exit signals to gracefully shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exit
		log.Println("closing postgres db connection")
		if err := sqlDB.Close(); err != nil {
			log.Println(fmt.Errorf("failed to close postgres db connection: %w", err))
		} else {
			log.Println("postgres db connection writer closed")
		}
	}()
	return &Postgres{sqlDB}, nil
}
