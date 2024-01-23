package repository

import "Quiklink_BE/internal/models"

type DatabaseRepo interface {
	RegisterUser(email, password, firstName, lastName, phone string) (models.User, error)
	LoginUser(email, password string) (models.User, error)
	AccessTokenLogin(token string) (models.User, error)
	InsertIntoShortUrlMap(shortURL, longURL string) error
	GetLongUrlFromShort(shortURL string) (string, error)
}
