package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {

	// use command-line flag to specify the network address
	// the default value is ":4000"
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// use the slog.NEW() function to create a new logger
	// which writes messages to the standard output stream
	// which write to the standard out stream and uses the
	// the default settings.
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))

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
	err := http.ListenAndServe(*addr, app.routes())

	// use Error() method to log the error message returned by
	// http.ListenAndServe at error severity level
	// and then call os.Exit(1) to terminate the program
	logger.Error(err.Error())
	os.Exit(1)

}
