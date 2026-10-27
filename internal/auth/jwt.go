package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


const SessionTTL = 24 * time.Hour

var CookieMaxAge = int(SessionTTL.Seconds())

var ErrInvalidToken = errors.New("invalid token")

// Claims son los claims custom + los registrados del JWT.
type Claims struct {
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"sub"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

// Sign emite un JWT HS256 firmado con el secreto.
func (m *JWTManager) Sign(tenantID, userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := Claims{
		TenantID: tenantID,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
			ID:        uuid.NewString(), // jti
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(m.secret)
}

// Verify valida firma, alg y expiración, y devuelve los claims.
func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		// Defensa contra alg-confusion: solo HS256 es aceptado.
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !tok.Valid {
		return nil, ErrInvalidToken
	}
	// Defensa en profundidad: un token sin tenant_id no es válido.
	if claims.TenantID == uuid.Nil || claims.UserID == uuid.Nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
