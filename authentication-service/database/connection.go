package database

import (

	"log"
	"os"
	"time"

	
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

// opens db connection
func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Cannot Open db ", err)
		return nil, err
	}
	err = db.AutoMigrate(&User{})
	if err !=nil{
		log.Println("Error in Migrations")
	}
	log.Println("migrations successful")
	//verifies if a connection to the database is still alive, establishing a connection if necessary.
	sqlDB, err := db.DB()
	if err !=nil {
		log.Println("error in getting sql")
	}
	
	err = sqlDB.Ping()
	if err != nil {
		log.Println("Error in connection", err)
		return nil, err
	}
	return db, nil
}

var counts int64

// connects to the database properly
func ConnectToDB() *gorm.DB {
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
