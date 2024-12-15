package auth

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type EmailProvider interface {
	SendEmail(content string) error
}

type AuthService struct {
	db      *sql.DB
	queries *database.Queries
	emailer EmailProvider
}

func NewAuthService(db *sql.DB, queries *database.Queries, emailer EmailProvider) *AuthService {
	return &AuthService{
		db:      db,
		queries: queries,
		emailer: emailer,
	}
}

func (as *AuthService) Register(ctx context.Context, email string, passowrd string) (*uuid.UUID, error) {

	tx, err := as.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("auth error begingin transaction: %w", err)
	}
	defer tx.Rollback()

	txq := as.queries.WithTx(tx)

	if !isValidCredentials(email, passowrd) {
		return nil, fmt.Errorf("auth error: invalid credintals")
	}

	exist, err := doesUserExistWithEmail(ctx, txq, email)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, fmt.Errorf("auth error: user already exist with email %v", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passowrd), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("auth error: genrating hashed password returned an error %v", err)
	}

	createUserParams := database.CreateUserParams{
		Email:    email,
		Password: string(hashedPassword),
	}
	user, err := txq.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, fmt.Errorf("auth error: creating new user in db %w", err)
	}

	// create a new verify email request
	token, err := createVerifyEmailRequest(ctx, txq, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating verfication request: %w", err)
	}
	// send verfity email
	verifiyEmailMSG := fmt.Sprintf("hello user, please verify your email with https://nuha.com/verify_email/%s", token)
	err = as.emailer.SendEmail(verifiyEmailMSG)
	if err != nil {
		return nil, fmt.Errorf("auth error sending email verfication: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("auth error commiting transaction: %w", err)
	}

	return &user.ID, nil
}

func (as *AuthService) VerfiyUser(ctx context.Context, email string, token string) error {

	verReq, err := getVerficationRequest(ctx, as.queries, token)
	if err != nil {
		return fmt.Errorf("auth error: error getting verfication token: %w", err)
	}

	if verReq.ExpiresAt.Compare(time.Now()) == -1 {
		return fmt.Errorf("auth error: verfication token has expired")
	}

	as.queries.SetUserVerified(ctx, email)
	return nil
}

func isValidCredentials(email string, pw string) bool {
	// check the email is valid email for registration
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isValidEmail, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return false
	}
	if !isValidEmail {
		return false
	}
	if len([]byte(pw)) > 72 {
		return false
	}

	// check password with regex (8 chars)
	validPasswordRegex := `^[A-Za-z0-9]{8,}$`

	isValidPassword, err := regexp.MatchString(validPasswordRegex, pw)
	if err != nil {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pw)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pw)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pw)

	return isValidPassword && hasUpper && hasLower && hasNumber
}
