package repository

import "Quiklink_BE/internal/models"

type DatabaseRepo interface {
	RegisterUser(email, password, firstName, lastName, phone string) (models.User, error)
	LoginUser(email, password string) (models.User, error)
	InsertIntoShortUrlMap(shortKey, longURL string) error
	GetLongUrlFromShort(shortKey string) (string, error)
}
