package http

import (
	"net/http"
	"rankit/errors"

	"github.com/alexedwards/scs/v2"
)

func Auth(sm *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userID := sm.GetString(r.Context(), "message")
			if userID != "" {
				RespondError(w, &errors.Error{
					Kind:    errors.Unauthorized,
					Message: "auth required",
				})
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
