package repositories

import (
	"context"
	"database/sql"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/google/uuid"
)

type UserRespository struct {
	db        *sql.DB
	dbQuereis *database.Queries
}

func NewUserRespository(db *sql.DB, dbQuereis *database.Queries) *UserRespository {
	return &UserRespository{
		db:        db,
		dbQuereis: dbQuereis,
	}
}

func (ur *UserRespository) DoesUserExistWithEmail(ctx context.Context, email string) (bool, error) {
	_, err := ur.dbQuereis.GetUserByEmail(ctx, email)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return true, nil
}

func (ur *UserRespository) StoreNewUserData(ctx context.Context, id uuid.UUID, firstName string, lastName string) (*models.User, error) {

	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	txq := ur.dbQuereis.WithTx(tx)

	userDB, err := txq.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	createUserDataParams := database.CreateUserDataParams{
		ID:        userDB.ID,
		FirstName: firstName,
		LastName:  lastName,
	}
	userDataDB, err := txq.CreateUserData(ctx, createUserDataParams)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        userDB.ID,
		FirstName: userDataDB.FirstName,
		LastName:  userDataDB.LastName,
		Email:     userDB.Email,
	}, nil

}

func (ur *UserRespository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	userDB, err := ur.dbQuereis.GetUserDataByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        userDB.ID,
		FirstName: userDB.FirstName,
		LastName:  userDB.LastName,
		Email:     userDB.LastName,
	}, nil
}
