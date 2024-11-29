package factories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserFactory struct {
	db  *database.Queries
	ctx context.Context
}

func NewUserFactory(ctx context.Context, db *database.Queries) *UserFactory {
	return &UserFactory{
		db:  db,
		ctx: ctx,
	}
}

func (uf *UserFactory) DoesUserExistWithEmail(email string) (bool, error) {
	_, err := uf.db.GetUserByEmail(uf.ctx, email)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return true, nil
}

func (uf *UserFactory) CreateNewUser(name string, email string, password string) (*models.User, error) {

	// TO-DO: hash the password here
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error genrating hashed password")
	}

	newUserParams := database.CreateUserParams{
		Name:     name,
		Email:    strings.ToLower(email),
		Password: string(hashedPassword),
	}
	user, err := uf.db.CreateUser(uf.ctx, newUserParams)
	if err != nil {
		return nil, err
	}

	return models.UserFromDBObject(&user), nil
}

func (uf *UserFactory) GetUserByEmail(email string) (*models.User, error) {

	email = strings.ToLower(email)
	user, err := uf.db.GetUserByEmail(uf.ctx, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return models.UserFromDBObject(&user), nil
}
