package middleware

import (
	"bitbucket.org/evaly/go-boilerplate/api/response"
	"net/http"
)

func AppKeyChecker(appKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Application-Key")
			if len(key) == 0 {
				response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "'Application-Key' required", "")
				return
			}
			if key != appKey {
				response.ServeJSON(w, http.StatusUnauthorized, nil, nil, "invalid 'Application-Key'", "")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
