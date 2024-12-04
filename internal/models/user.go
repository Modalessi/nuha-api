package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateUser(name string, email string, password string) (*User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error genrating hashed password")
	}

	return &User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}, nil
}

func UserFromDBObject(u *database.User) *User {
	return &User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func (u *User) JSON() []byte {
	data, err := json.Marshal(u)
	utils.Assert(err, "error converting user object to json")
	return data
}
