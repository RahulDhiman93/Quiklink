package dbrepo

import (
	"Quiklink_BE/internal/models"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// RegisterUser Register the user
func (m *postgresDBRepo) RegisterUser(email, password, firstName, lastName, phone string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	var existingEmail string
	err := m.DB.QueryRowContext(ctx, "SELECT email FROM users WHERE email = $1", email).Scan(&existingEmail)
	if err == nil {
		return user, fmt.Errorf("email %s is already registered", email)
	} else if !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error checking existing email:", err)
		return user, err
	}

	query := `insert into users (access_token, first_name, last_name, email, password, phone,created_at,updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating bcrypt hash:", err)
		return user, err
	}

	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return user, err
	}
	token := base64.URLEncoding.EncodeToString(randomBytes)

	_, err = m.DB.ExecContext(ctx, query, token, firstName, lastName, email, hashedPassword, phone, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return user, err
	}

	query = `SELECT id, access_token, first_name, last_name, email, password, phone, access_level FROM users WHERE email = $1`

	err = m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.AccessToken,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.AccessLevel,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

// LoginUser Login the user
func (m *postgresDBRepo) LoginUser(email, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	query := `SELECT id, access_token, first_name, last_name, email, password, phone, access_level FROM users WHERE email = $1`

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.AccessToken,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.AccessLevel,
	)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return user, errors.New("incorrect password")
	}
	return user, nil
}

func (m *postgresDBRepo) InsertIntoShortUrlMap(shortURL, longURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if exists, err := m.shortURLExists(ctx, shortURL); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("shortURL %s already exists", shortURL)
	}

	query := `insert into short_url_map (short_url, long_url, created_at, updated_at)
              values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, query,
		shortURL,
		longURL,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) shortURLExists(ctx context.Context, shortURL string) (bool, error) {
	query := "select exists (select 1 from short_url_map where short_url = $1)"

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, shortURL).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetLongUrlFromShort gets a long url from short
func (m *postgresDBRepo) GetLongUrlFromShort(shortURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var longURL string

	query := `select long_url from short_url_map where short_url = $1`

	row := m.DB.QueryRowContext(ctx, query, shortURL)
	err := row.Scan(
		&longURL,
	)

	if err != nil {
		return "", err
	}

	return longURL, nil
}
