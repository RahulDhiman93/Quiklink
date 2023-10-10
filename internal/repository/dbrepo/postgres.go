package dbrepo

import (
	"context"
	"time"
)

func (m *postgresDBRepo) InsertIntoShortUrlMap(shortURL, longURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
