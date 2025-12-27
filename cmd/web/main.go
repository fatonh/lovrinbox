package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
}

func main() {

	// use command-line flag to specify the network address
	// the default value is ":4000"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// define a new command-line flag for the MYSQL DSN string
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true",
		"MySQL data source name")

	flag.Parse()

	// use the slog.NEW() function to create a new logger
	// which writes messages to the standard output stream
	// which write to the standard out stream and uses the
	// the default settings.
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))

	//to keep the main() function tidy i've put the code for crateing
	// a connection pool into the separet openDB function below.
	// We pass openDB the DSN string as a parameter form
	// command-line flag
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)

	}

	// we also defer a call to db.Close(), so that the
	// connection pool is closed before the main() function exits
	defer db.Close()

	// Initialize a new instance of application containing
	// the dependencies for our application struct, containing
	// the dependencies (for noew, just the struct logger)
	app := &application{
		logger: logger,
	}

	// use the Info() method to log the starting server message
	// at info severity level
	logger.Info("Starting server", "addr", *addr)

	// use the htt.ListenAndServe function to start
	// a new web server. We pass in two parameters: The TCP netweokd address
	// to listen on (in this case ":4000") and the servemux
	// we just created. If http.ListenAndServe returns an error, we log it and exit
	// we use log.Fatal function
	// to log the error message and terminate the program.
	// Note that any error returned by http.ListenAndServe is
	// allways non-nil.
	err = http.ListenAndServe(*addr, app.routes())

	// use Error() method to log the error message returned by
	// http.ListenAndServe at error severity level
	// and then call os.Exit(1) to terminate the program
	logger.Error(err.Error())
	os.Exit(1)

}

// the openDB() function wraps sql.Open() and
// returns a sql.DB connection pool for a given DSN string
func openDB(dsn string) (*sql.DB, error) {
	// use sql.Open to create an empty connection pool
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// use the db.Ping() method to establish a new connection
	// to the database, which verifies that the DSN is valid
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
