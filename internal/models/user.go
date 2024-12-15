package models

import (
	"encoding/json"

	"github.com/Modalessi/nuha-api/internal/utils"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

func (u *User) JSON() []byte {
	data, err := json.Marshal(u)
	utils.Assert(err, "error converting user object to json")
	return data
}
