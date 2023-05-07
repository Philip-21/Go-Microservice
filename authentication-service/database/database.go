package database

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (u *User) CreateUser(firstname string, lastname string, email string, password string) (*User, error) {
	query :=
		`INSERT INTO users
		(first_name, last_name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id, first_name, last_name, email, password, created_at, updated_at`

	var user User
	//executing raw queries with gorm
	if err := db.Raw(query, firstname, lastname, email, password, time.Now(), time.Now()).Scan(&user).Error; err != nil {
		return nil, err
	}

	log.Println("User Created")
	return &user, nil
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	var users []*User
	if err := db.Order("last_name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByEmail returns one user by email
// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id int) (*User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() error {
	return db.Model(&User{}).Where("id = ?", u.ID).Updates(User{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Active:    u.Active,
		UpdatedAt: time.Now(),
	}).Error
}

// Delete deletes one user from the database, by User.ID
func (u *User) Delete() error {
	return db.Delete(&User{}, u.ID).Error
}

// DeleteByID deletes one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	return db.Delete(&User{}, id).Error
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	// create a new User object with the hashed password
	newUser := &User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  string(hashedPassword),
		Active:    user.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// insert the new user into the database using a raw SQL query
	result := db.Raw("INSERT INTO users (email, first_name, last_name, password, user_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id",
		newUser.Email, newUser.FirstName, newUser.LastName, newUser.Password, newUser.Active, newUser.CreatedAt, newUser.UpdatedAt).Row()

	// get the new user's ID
	var newID int
	if err := result.Scan(&newID); err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// update the user's password using a raw SQL query
	stmt := "UPDATE users SET password = ? WHERE id = ?"
	if err := db.Exec(stmt, string(hashedPassword), u.ID).Error; err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
