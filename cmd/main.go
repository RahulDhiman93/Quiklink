package main

import (
	"Quiklink_BE/internal/driver"
	"Quiklink_BE/internal/handlers"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const portNum = ":8080"

func main() {
	db, err := run()
	if err != nil {
		log.Println("Error in initial RUN")
		log.Fatal(err)
	}
	defer db.SQL.Close()

	srv := &http.Server{
		Addr:    portNum,
		Handler: routes(),
	}

	fmt.Println("Application listening on ", portNum)
	err = srv.ListenAndServe()
	if err != nil {
		log.Println("Error in LISTEN and SERVE")
	}
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	//read flags
	dbName := flag.String("dbname", "", "Database name")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database pass")
	dbPort := flag.String("dbport", "5432 ", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	//Connect to DB
	log.Println("<<<-- Connecting to DB")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	log.Println(connectionString)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to DB, dying!!!...")
		return nil, err
	}
	log.Println("Connected to DB -->>>")

	repo := handlers.NewRepo(db)
	handlers.FreshHandlers(repo)

	return db, nil
}
