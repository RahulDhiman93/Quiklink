package main

import (
	"Quiklink_BE/internal/config"
	"Quiklink_BE/internal/driver"
	"Quiklink_BE/internal/handlers"
	"Quiklink_BE/internal/helpers"
	"Quiklink_BE/internal/models"
	"Quiklink_BE/internal/render"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const portNum = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Println("Error in initial RUN")
		log.Fatal(err)
	}
	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNum,
		Handler: routes(&app),
	}

	fmt.Println("Application listening on ", portNum)
	err = srv.ListenAndServe()
	if err != nil {
		log.Println("Error in LISTEN and SERVE")
	}
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	gob.Register(models.User{})
	gob.Register(models.TemplateData{})
	gob.Register(models.AuthRequestBody{})
	gob.Register(map[string]int{})

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	//read flags
	inProduction, _ := strconv.ParseBool(os.Getenv("IN_PRODUCTION"))
	useCache, _ := strconv.ParseBool(os.Getenv("USE_CACHE"))
	dbName := os.Getenv("DBNAME")
	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbPort := os.Getenv("DBPORT")
	dbSSL := os.Getenv("DBSSL")

	if dbName == "" || dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	//change to true for Production
	app.InProduction = inProduction
	app.UseCache = useCache

	//Info Log setup
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	//Error Log setup
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//Connect to DB
	log.Println("<<<-- Connecting to DB")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL)
	log.Println(connectionString)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to DB, dying!!!...")
	}
	log.Println("Connected to DB -->>>")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepository(&app, db)
	handlers.FreshHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
