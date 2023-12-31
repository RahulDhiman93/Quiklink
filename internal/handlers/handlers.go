package handlers

import (
	"Quiklink_BE/internal/driver"
	"Quiklink_BE/internal/models"
	"Quiklink_BE/internal/render"
	"Quiklink_BE/internal/repository"
	"Quiklink_BE/internal/repository/dbrepo"
	"encoding/json"
	"github.com/go-chi/chi"
	qrcode "github.com/skip2/go-qrcode"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var Repo *Repository

type Repository struct {
	DB repository.DatabaseRepo
}

// NewRepo creates a new repo
func NewRepo(db *driver.DB) *Repository {
	return &Repository{
		DB: dbrepo.NewPostgresRepo(db.SQL),
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

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = render.TemplateRenderer(w, r, "home.page.tmpl", &models.TemplateData{})
}

// ShortenURL generates a short key for a given URL and stores it in the map.
func (m *Repository) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var request struct {
		LongURL string `json:"long_url"`
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

	shortKey := randomString(8)

	err := m.DB.InsertIntoShortUrlMap(shortKey, request.LongURL)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	//shortUrl := "http://localhost:8080/" + shortKey //DEV
	//shortUrl := "https://ec2-18-144-176-134.us-west-1.compute.amazonaws.com/" + shortKey //REVERSE DNS
	shortUrl := "https://quiklink.site/" + shortKey //PROD

	code, err := generateQRCode(shortUrl)
	if err != nil {
		log.Println(err)
	}

	response := jsonResponse{
		OK:       true,
		Message:  "Short URL created",
		LongURL:  request.LongURL,
		ShortURL: shortUrl,
		QRCode:   code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Redirect redirects the short URL to the original long URL.
func (m *Repository) Redirect(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")
	longURL, err := m.DB.GetLongUrlFromShort(shortKey)

	if err != nil || longURL == "" {
		stringMap := make(map[string]string)
		stringMap["url_not_found"] = "url_not_found"
		_ = render.TemplateRenderer(w, r, "home.page.tmpl", &models.TemplateData{
			StringMap: stringMap,
		})
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func generateQRCode(url string) ([]byte, error) {
	data := url
	code, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		return nil, err
	}

	return code, nil
}

// RandomString Generates a random string
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}
