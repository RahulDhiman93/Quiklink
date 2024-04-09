package handlers

import (
	"Quiklink_BE/internal/config"
	"Quiklink_BE/internal/driver"
	"Quiklink_BE/internal/models"
	"Quiklink_BE/internal/render"
	"Quiklink_BE/internal/repository"
	"Quiklink_BE/internal/repository/dbrepo"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepository creates a new repo
func NewRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new test repo
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// FreshHandlers sets the repository for the handlers
func FreshHandlers(r *Repository) {
	Repo = r
}

type jsonResponse struct {
	OK       bool   `json:"ok"`
	Message  string `json:"message"`
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
	QRCode   []byte `json:"qrcode"`
}

type authResponse struct {
	OK      bool                   `json:"ok"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplateRenderer(w, r, "home.page.tmpl", &models.TemplateData{})
}

// RegisterUser creates a new user
func (m *Repository) RegisterUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		internalServerError(w, err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	phone := r.Form.Get("phone")

	if email == "" || password == "" || firstName == "" || lastName == "" {
		internalServerError(w, fmt.Errorf("please add all attributes"))
		return
	}

	user, err := m.DB.RegisterUser(email, password, firstName, lastName, phone)
	if err != nil {
		internalServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", user.Id)
	m.App.Session.Put(r.Context(), "flash", "Registration successful!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LoginUser creates a new user
func (m *Repository) LoginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		internalServerError(w, err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" || password == "" {
		internalServerError(w, fmt.Errorf("please add all attributes"))
		return
	}

	user, err := m.DB.LoginUser(email, password)
	if err != nil {
		internalServerError(w, err)
		return
	}

	respData := make(map[string]interface{})
	respData["user"] = user
	response := authResponse{
		OK:      true,
		Message: "user authenticated successfully",
		Data:    respData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ShortenURL generates a short key for a given URL and stores it in the map.
func (m *Repository) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var request struct {
		LongURL  string `json:"long_url"`
		UrlAlias string `json:"url_alias"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if strings.Contains(request.LongURL, "quiklink.site") {
		resp := jsonResponse{
			OK:      false,
			Message: "Quiklink URL not allowed",
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
		}
		return
	}

	shortKey := request.UrlAlias
	if request.UrlAlias == "" {
		shortKey = randomString(8)
	}
	log.Println("SHORTKEY -->", shortKey)

	baseUrl := "https://quiklink.site/" //PROD
	//baseUrl := "http://localhost:8080/" //DEV

	errChan := make(chan error)
	defer close(errChan)

	go func() {
		errChan <- m.DB.InsertIntoShortUrlMap(shortKey, request.LongURL)
	}()

	qrCodeChan := make(chan []byte)
	defer close(qrCodeChan)

	go func() {
		code, err := generateQRCode(baseUrl + shortKey)
		if err != nil {
			log.Println(err)
			qrCodeChan <- []byte{}
			return
		}
		qrCodeChan <- code
	}()

	dbErr := <-errChan
	qrCode := <-qrCodeChan

	if dbErr != nil {
		respondWithError(w, dbErr.Error(), http.StatusBadRequest)
		return
	}

	response := jsonResponse{
		OK:       true,
		Message:  "Short URL created",
		LongURL:  request.LongURL,
		ShortURL: baseUrl + shortKey,
		QRCode:   qrCode,
	}

	respondWithJSON(w, response, http.StatusOK)
}

// Redirect redirects the short URL to the original long URL.
func (m *Repository) Redirect(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")

	longUrlChan, errChan := make(chan string), make(chan error)
	defer close(longUrlChan)
	defer close(errChan)

	go func() {
		longURL, err := m.DB.GetLongUrlFromShort(shortKey)
		longUrlChan <- longURL
		errChan <- err
	}()

	longUrlDB, dbErr := <-longUrlChan, <-errChan

	if dbErr != nil || longUrlDB == "" {
		stringMap := make(map[string]string)
		stringMap["url_not_found"] = "url_not_found"
		_ = render.TemplateRenderer(w, r, "home.page.tmpl", &models.TemplateData{
			StringMap: stringMap,
		})
		return
	}

	http.Redirect(w, r, longUrlDB, http.StatusSeeOther)
}
