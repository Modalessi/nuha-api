package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/models"
)

type UserRespository struct {
	db  *database.Queries
	ctx context.Context
}

func NewUserRespository(ctx context.Context, db *database.Queries) *UserRespository {
	return &UserRespository{
		db:  db,
		ctx: ctx,
	}
}

func (ur *UserRespository) DoesUserExistWithEmail(email string) (bool, error) {
	_, err := ur.db.GetUserByEmail(ur.ctx, email)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return true, nil
}

func (uf *UserRespository) StoreNewUser(u *models.User) (*database.User, error) {

	newUserParams := database.CreateUserParams{
		Name:     u.Name,
		Email:    strings.ToLower(u.Name),
		Password: u.Password,
	}

	user, err := uf.db.CreateUser(uf.ctx, newUserParams)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRespository) GetUserByEmail(email string) (*database.User, error) {

	email = strings.ToLower(email)
	user, err := ur.db.GetUserByEmail(ur.ctx, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
