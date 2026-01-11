package main

import (
	"fmt"
	"net/http"
)

// middleware pattern.
// because we want this middleware to act on every request this is received
// the chain of controll looks like this:
// commonHeaders -> servemux -> application handler
func commonHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: this is a split across multiple lines for readability.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "GO")

		next.ServeHTTP(w, r)
	})

}

// logRequest is a middleware which logs basic information about each
// HTTP request received.
// logRequest -|> commonHeaders -|> servemux -|> application handler
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto,
			"method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a deferred function which recovers from a panic
		defer func() {
			// use the built-in recover() function to check if the panic occurred.
			// if a panic did happen, recover() will the return the panic value. if
			// a panic didn't happen, it will return nil.

			pv := recover()

			// if a panic did happen ...
			if pv != nil {
				// set a "connection: close" header on the response
				w.Header().Set("Connection", "close")

				// call the app.serverError() helper method to log the error
				// and send the 500 Internal Server Error response to the user
				app.serverError(w, r, fmt.Errorf("%v", pv))
			}
		}()

		// call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
