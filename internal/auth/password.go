package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword genera un bcrypt con cost por defecto (10).
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword devuelve true si password corresponde a hash.
// Nunca propaga el error de bcrypt: hacia fuera solo importa sí/no.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
