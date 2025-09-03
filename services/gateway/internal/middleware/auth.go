package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type AuthConfig struct{ JWKSURL string }

func JWTAuth(cfg AuthConfig) func(http.Handler) http.Handler {
	var set jwk.Set
	if cfg.JWKSURL != "" {
		// Load once on startup; in prod, refresh periodically
		if s, err := jwk.Fetch(context.Background(), cfg.JWKSURL); err == nil {
			set = s
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ah := r.Header.Get("Authorization")
			if ah == "" || !strings.HasPrefix(ah, "Bearer ") {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			tokstr := strings.TrimPrefix(ah, "Bearer ")
			var opt []jwt.ParseOption
			if set != nil {
				opt = append(opt, jwt.WithKeySet(set))
			}
			if _, err := jwt.Parse([]byte(tokstr), opt...); err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
