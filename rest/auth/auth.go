// Package auth provides HTTP authentication middleware for the REST service.
//
// It mirrors the bearer-token middleware in
// rest_api_clientandserver/server_http_go/main.go: a request must carry the
// header "Authorization: Bearer <token>" or it is rejected with 401
// Unauthorized. Because Auth wraps an http.Handler (the whole router), every
// route behind it is protected in one place — the wiring layer in main.go —
// instead of repeating the check inside each handler.
package auth

import "net/http"

// expectedToken is the bearer credential the middleware requires.
// In a real service this would come from config or a secret store, not a
// hard-coded constant — it is kept literal here to match the reference example.
const expectedToken = "Bearer abcxyz"

// Auth wraps next and lets the request through only when the Authorization
// header exactly matches the expected bearer token; otherwise it short-circuits
// with 401 Unauthorized and next is never called.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != expectedToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
