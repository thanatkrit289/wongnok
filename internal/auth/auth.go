package auth

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// Auth service
type Auth struct {
	db *sql.DB
	repo repository
}

type repository interface {
	InsertUser(ctx context.Context, db *sql.DB, username, password string) (userID int64, err error)
}

// New creates new auth service
func New(db *sql.DB) *Auth {
	return &Auth{db, repo{}}
}

var reUsername = regexp.MustCompile(`^[a-z0-9]*$`)

// Sign up
func (svc *Auth) SignUp(ctx context.Context, username, password string) (userID int64, err error) {
	// normalize data
	username = strings.ToLower(username)
	username = strings.TrimSpace(username)

	// validate
	if username == "" {
		return 0, fmt.Errorf("username required")
	}
	if len(username) > 20 {
		return 0, fmt.Errorf("username too long")
	}
	if !reUsername.MatchString(username) {
		return 0, fmt.Errorf("invalide format username")
	}
	if password == "" {
		return 0, fmt.Errorf("invalide format username")
	}
	if len(password) < 6 {
		return 0, fmt.Errorf("password too shart")
	}
	if len(password) > 64 {
		return 0, fmt.Errorf("password too long")
	}

	// hash password
	hashedPass := HashPassword(password)

	userID, err = svc.repo.InsertUser(ctx, svc.db, username, hashedPass)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// SignIn sign in user
func (svc *Auth) SignIn(ctx context.Context, username, password string) (token string, err error) {
	username = strings.ToLower(username)
	username = strings.TrimSpace(username)

	if username == "" {
		return "", fmt.Errorf("username required")
	}
	if len(username) > 20 {
		return "", fmt.Errorf("username too long")
	}
	if password == "" {
		return "", fmt.Errorf("password required")
	}
	if len(password) > 64 {
		return "", fmt.Errorf("password too long")
	}

	var (
		userID       string
		userPassword string
	)
	err = svc.db.QueryRowContext(ctx, `
		select 
			id, password
		from users
		where username = $1
	`, username).Scan(&userID, &userPassword)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return "", err
	}

	if !compareHashAndPassword(userPassword, password) {
		return "", fmt.Errorf("invalid credentials")
	}

	// generate token
	token = generateToken()

	// insert token to data base
	svc.db.ExecContext(ctx, `
		insert into auth_tokens
			(id, user_id)
		values
			($1, $2)
	`, token, userID)
	if err != nil {
		return "", err
	}

	// return token
	return token, nil
}

// SignOut sign out user
func (svc *Auth) SignOut(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("token required")
	}

	_, err := svc.db.ExecContext(ctx, `
		delete from auth_tokens
		where id = $1
	`, token)

	return err
}
