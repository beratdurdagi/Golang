package Types

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// USER Struct for Login and event owners
type User struct {
	ID          int       `json:"_id"`
	Name        string    `json:"first_name"`
	Surname     string    `json:"last_name"`
	Email       string    `json:"email"`
	EncPassword string    `json:"enc_password"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
	Deleted_At  time.Time `json:"deleted_at"`
}

// Constructor
func NewAccount(fistName string, lastName string, email string, password string) (*User, error) {

	return &User{

		Name:        fistName,
		Surname:     lastName,
		EncPassword: password,
		Email:       email,
		Created_At:  time.Now().UTC(),
		Updated_At:  time.Now().UTC(),
	}, nil

}

// Validate Func to  Validate Login Req password and hashed password
func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncPassword), []byte(password)) == nil
}

// Login Response for given
type LoginResponse struct {
	UserId int `json:"id"`
	Email  string
	Token  string `json:"token"`
}

// Login Req
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//SignUp Request

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Response
type Response struct {
	Message string `json:"message"`
	Error   error  `json:"status"`
}
