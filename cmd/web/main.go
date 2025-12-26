package main

import (
	"log"
	"net/http"
)

func main() {
	// use the http.NewServeMux to create a new ServeMux
	// then register the home function as the handler for the "/" route
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() method to register the file server as the handler for
	// all URL paths starting with "/static/". To do this, we use the
	// http.StripPrefix() function to modify the request URL path before the
	// file server sees it. If we didn't do this, the file server would look for
	// files with paths like "./ui/static/static/css/main.css", which obviously
	// don't exist. By using http.StripPrefix(), the file server will see a path
	// like "./ui/static/css/main.css", which is correct.
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/craete", snippetCreatePost)

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
