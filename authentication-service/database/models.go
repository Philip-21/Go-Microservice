package database

import (
	"time"

	"gorm.io/gorm"
)

const dbTimeout = time.Second * 3

var db *gorm.DB

// User is the structure which holds one user from the database.
type User struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Email     string `gorm:"unique" json:"email"`
	FirstName string `gorm:"column:first_name" json:"first_name,omitempty"`
	LastName  string `gorm:"column:last_name" json:"last_name,omitempty"`
	Password  string `gorm:"-" json:"-"`
	Active    int `gorm:"default:1" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User //using User 2 define the object to be used in the db reciver func
}

// New is the function used to create an instance of the data package by connecting to the db.
// It returns the type Model, which embeds all the types we want to be available to our application.
func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}
