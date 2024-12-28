package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/golang-jwt/jwt/v5"
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
	config  AuthServiceConfig
}

type AuthServiceConfig struct {
	JWTSecretKey              string
	TokensExpirationsDuration time.Duration
}

func NewAuthService(db *sql.DB, queries *database.Queries, emailer EmailProvider, config AuthServiceConfig) *AuthService {
	return &AuthService{
		db:      db,
		queries: queries,
		emailer: emailer,
		config:  config,
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

func (as *AuthService) VerifyUser(ctx context.Context, token string) error {
	tx, err := as.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("auth error: starting transaction: %w", err)
	}
	defer tx.Rollback()

	txq := as.queries.WithTx(tx)

	verReq, err := getVerficationRequest(ctx, txq, token)
	if err != nil {
		return fmt.Errorf("auth error: error getting verification token: %w", err)
	}

	if verReq.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("auth error: verification token has expired")
	}

	user, err := txq.GetUserByID(ctx, verReq.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Verified {
		return fmt.Errorf("user is already verified")
	}

	_, err = txq.SetUserVerified(ctx, verReq.UserID)
	if err != nil {
		return fmt.Errorf("auth error setting user verified: %w", err)
	}

	_, err = txq.DelteVerficationToken(ctx, token)
	if err != nil {
		return fmt.Errorf("auth error deleting verification token: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("auth error committing changes: %w", err)
	}

	return nil
}

func (as *AuthService) Login(ctx context.Context, email string, password string) (string, error) {

	user, err := as.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("auth error: user does not exist")
		}
		return "", err
	}

	// check password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("auth error: wrong credintals")
	}

	// check user is verfied
	if !user.Verified {
		return "", fmt.Errorf("auth error: user is not verfied")
	}

	// create jwt token here

	token, err := NewJWTTokenWithClaims(user.Email, time.Now().Add(as.config.TokensExpirationsDuration), as.config.JWTSecretKey)
	if err != nil {
		return "", fmt.Errorf("auth error: creating jwt token: %w", err)
	}

	// password is correct so create session, and return the created tokn
	createSessionParams := database.CreateUserSessionParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(as.config.TokensExpirationsDuration),
	}
	_, err = as.queries.CreateUserSession(ctx, createSessionParams)
	if err != nil {
		return "", fmt.Errorf("auth error: creating session: %v", err)
	}

	return token, nil
}

func (as *AuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	// first verify token
	jwtToken, err := VerfiyToken(token, as.config.JWTSecretKey)
	if err != nil {
		return "", fmt.Errorf("auth error: invalid jwt token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("auth error: error getting claims from token")
	}

	userEmail, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("auth error: sub key was not found in token")
	}

	// get session and check it is not revoked
	session, err := as.queries.GetSession(ctx, token)
	if err != nil {
		return "", fmt.Errorf("auth error: getting session from db %v", err)
	}

	// check session is not invalid
	if session.Revoked {
		return "", fmt.Errorf("auth error: session revoked")
	}

	return userEmail, nil
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
