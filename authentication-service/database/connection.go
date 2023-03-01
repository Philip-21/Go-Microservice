package database

import (
	"context"
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
		log.Println("Cannot Open db ", err)
		return nil, err
	}
	//verifies if a connection to the database is still alive, establishing a connection if necessary.
	err = db.Ping()
	if err != nil {
		log.Println("Error in connection", err)
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

func SeedDB(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id Serial Primary key,
		email character varying(225) NOT NULL,
		first_name character varying(225) NOT NULL,
		last_name character varying(225) NOT NULL, 
		password character varying(225) NOT NULL,
		user_active int default 0 NOT NULL,
		created_at timestamp(6) NOT NULL,
		updated_at timestamp(6) NOT NULL
	) `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating table ", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Errors making rows affected", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("Error %s when pinging database\n", err)
	}
	log.Printf("Users table created table :%d", rows)
	log.Println("Users table created")
	return nil

}
