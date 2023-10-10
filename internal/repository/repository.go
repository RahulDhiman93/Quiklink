package repository

type DatabaseRepo interface {
	InsertIntoShortUrlMap(shortURL, longURL string) error
	GetLongUrlFromShort(shortURL string) (string, error)
}
