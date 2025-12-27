package main

import "net/http"

func (app *application) routes() *http.ServeMux {

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

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/craete", app.snippetCreatePost)

	return mux
}
