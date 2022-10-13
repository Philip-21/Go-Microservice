package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// opens db connection
func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	//verifies if a connection to the database is still alive, establishing a connection if necessary.
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

var counts int64

// connects to the database properly
func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN") //DSN is gotten from the docker compose yml file in auth service
	//an infinite for loop to connect to the database
	for {
		connection, err := OpenDB(dsn)
		if err != nil {
			log.Println("Postgres not yet connected......")
			counts++ //add 1 to counts
		} else {
			log.Println("Connected to Postgres")
			return connection
		}
		//tryin to connect to database for 20 seconds
		//having 10 counts
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing of for two seconds.......")
		time.Sleep(2 * time.Second) //waiting for 2sec each time
		continue
	}
}
