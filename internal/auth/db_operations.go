package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/google/uuid"
)

func doesUserExistWithEmail(ctx context.Context, txq *database.Queries, email string) (bool, error) {
	_, err := txq.GetUserByEmail(ctx, email)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return true, nil
}

func createVerifyEmailRequest(ctx context.Context, txq *database.Queries, userID uuid.UUID) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("error generating random token: %w", err)
	}

	verficationParams := database.CreateVerificationRequestParams{
		UserID: userID,
		Token:  token,
	}
	verficationReq, err := txq.CreateVerificationRequest(ctx, verficationParams)
	if err != nil {
		return "", fmt.Errorf("error creating verfication request in db: %w", err)
	}

	return verficationReq.Token, nil
}

func getVerficationRequest(ctx context.Context, txq *database.Queries, token string) (*database.VerificationToken, error) {
	verReq, err := txq.GetVerficationToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &verReq, nil
}

func generateToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", token), nil
}
