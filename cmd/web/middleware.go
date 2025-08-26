package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack)
		defer func() {
			// use the builtin  recover function to check if there has been
			// a panic or not. If there has...
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if user is not authenticated, redirect to login page
		// and return from the middleware chain so that no sobsequent
		// handlers in the chain are executed
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// set the "Cache-Control" header so that pages require authentication aren't cached
		w.Header().Set("Cache-Control", "no-store")

		// call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: true,
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve authenticatedUserID from the session.
		// will return the 0 value for an int (0) if no aithenticatedUserID
		// value is found  -- in which case we call the next handler
		// in the chain as normal and return
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// chack if an user with the id exists in the database
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// is a matching user is found, create a new copy of the request
		// (with an isAuthenticatedContextKey value of true in the context)
		// ans assign it to r
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		// call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
