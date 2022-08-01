package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
	"weight.kenfan.org/internal/models"
)

// application struct to hold app-wide dependencies
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	weights       *models.WeightModel
	templateCache map[string]*template.Template
}

func main() {
	// Create cmd line flag in addr var
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "/Users/kenny/Documents/weight/weight.db", "SQLite data source")

	flag.Parse()
	// Logger for writing info messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		weights:       &models.WeightModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
