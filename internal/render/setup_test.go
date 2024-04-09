package render

import (
	"Quiklink_BE/internal/config"
	"Quiklink_BE/internal/helpers"
	"Quiklink_BE/internal/models"
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	//Session date types
	gob.Register(models.User{})
	gob.Register(models.TemplateData{})
	gob.Register(models.AuthRequestBody{})
	gob.Register(map[string]int{})

	//change to true for Production
	testApp.InProduction = false

	//Info Log setup
	infoLog := log.New(os.Stdout, "INDO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	//Error Log setup
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp
	helpers.NewHelpers(app)

	os.Exit(m.Run())
}
