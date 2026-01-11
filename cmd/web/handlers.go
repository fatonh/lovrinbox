package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fatonh/lovrinbox/internal/models"
)

// change the signature of the home handler function
// so it is defined as method against *application
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//panic("oops! something went wrong") // deliberate panic for testing.

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year)
	// and then add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper method.
	app.render(w, r, http.StatusOK, "home.tmpl",
		data)

	// // Initialize a slice containing the paths to the two files. It's important
	// // to note that the file containing our base template must be the *first*
	// // file in the slice.
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }

	// // Use the template.ParseFiles() function to read the template file into a
	// // template set. If there's an error, we log the detailed error message, use
	// // the http.Error() function to send an Internal Server Error response to the
	// // user, and then return from the handler so no subsequent code is executed.
	// ts, err := template.ParseFiles(files...)

	// // Create an instance f templateData struct holding the slice of snippets.
	// data := templateData{
	// 	Snippets: snippets,
	// }

	// if err != nil {
	// 	app.serverError(w, r, err)

	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// // Use the ExecuteTemplate() method to write the content of the "base"
	// // template as the response body.
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// w.Write([]byte("Welcome to the Home Page!"))
}

// Add a snippetView handler function.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// use the SnippetModel's Get() method to retrieve the
	// the data for a specific record based on it's ID.
	// If no matching record is found. return a 404 Not found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// use the new render helper.
	app.render(w, r, http.StatusOK, "view.tmpl",
		data)

	// // Initialize a slice containing the paths to the view.tmpl files.
	// // plus the base layout and the nav partials.
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }

	// // Parse the template files ...
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// // Create an instance of a templateData struct to holding the snippet data.
	// data := templateData{
	// 	Snippet: snippet,
	// }

	// // And then execute them. Notice how we are passing in the snippet
	// // data (a models.Snipet struct) as final parameter.
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)

	// }

	// // write the snippet data as plain-text HTTP response body
	// fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet lovrin..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id),
		http.StatusSeeOther)
}
