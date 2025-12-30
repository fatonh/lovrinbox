package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

func (app *application) render(w http.ResponseWriter, r *http.Request,
	status int, page string, data templateData) {

	// Retrieve the appropriate template set from cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper
	// method that we made earlier and return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the tempalte %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there is an error, call our serverError()
	// helper method and return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//If the template is written to the buffer successfully, we
	// can then write the contents of the buffer to the
	// http.ResponseWriter.
	// Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	// Note: this is another time where we pass our http.ResponseWriter to
	// a function that takes an io.Writer.
	buf.WriteTo(w)

	// // Execute the template set and write the response body. Again if there
	// // is any error we call the serverError() helper method and return.
	// err := ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}

// Create a newTemplateData() helper. which returns a templateData struct
// initialized with the current year. Note that we're not using the *http.Request
// parameter here at the moment, but we will use it later in the book.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}
