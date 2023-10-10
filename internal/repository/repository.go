package repository

type DatabaseRepo interface {
	InsertIntoShortUrlMap(shortURL, longURL string) error
}
