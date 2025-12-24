package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Home Page!"))
}

// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// use the http.NewServeMux to create a new ServeMux
	// then register the home function as the handler for the "/" route
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("starting server on :4000")

	// use the htt.ListenAndServe function to start
	// a new web server. We pass in two parameters: The TCP netweokd address
	// to listen on (in this case ":4000") and the servemux
	// we just created. If http.ListenAndServe returns an error, we log it and exit
	// we use log.Fatal function
	// to log the error message and terminate the program.
	// Note that any error returned by http.ListenAndServe is
	// allways non-nil.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
