package main

import (
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter,
	r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// use debug.Stack() to get the stack trace.
		// this returns a byte slice, which we convert to a string
		// so it's readable in the log entry.
		trace = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method,
		"uri", uri, "trace", trace)

	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)

}

func (app *application) clientError(w http.ResponseWriter,
	r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}
