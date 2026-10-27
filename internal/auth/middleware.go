package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const CookieName = "session"

type ctxKey int

const (
	ctxKeyTenant ctxKey = iota
	ctxKeyUser
)

// TenantFromContext devuelve el tenant inyectado por Middleware.
func TenantFromContext(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(ctxKeyTenant).(uuid.UUID)
	return v, ok
}

// UserFromContext devuelve el user_id inyectado por Middleware.
func UserFromContext(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(ctxKeyUser).(uuid.UUID)
	return v, ok
}

// Middleware extrae la cookie de sesión, verifica el JWT e inyecta tenant_id y user_id en el contexto del request.
func Middleware(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(CookieName)
			if err != nil {
				writeUnauthorized(w, "cookie de sesión requerida")
				return
			}
			claims, err := jwtManager.Verify(c.Value)
			if err != nil {
				writeUnauthorized(w, "token inválido o expirado")
				return
			}
			ctx := context.WithValue(r.Context(), ctxKeyTenant, claims.TenantID)
			ctx = context.WithValue(ctx, ctxKeyUser, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SetSessionCookie emite la cookie con los atributos de seguridad acordados.
func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   CookieMaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// writeUnauthorized emite el error JSON uniforme.
func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}
