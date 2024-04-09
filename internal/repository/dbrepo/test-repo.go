package dbrepo

import (
	"Quiklink_BE/internal/models"
)

func (m *testDBRepo) RegisterUser(email, password, firstName, lastName, phone string) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) LoginUser(email, password string) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) InsertIntoShortUrlMap(shortKey, longURL string) error {
	return nil
}

func (m *testDBRepo) GetLongUrlFromShort(shortKey string) (string, error) {
	s := ""
	return s, nil
}
