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
