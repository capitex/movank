package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// idempotencyHash calcula SHA-256 sobre la serialización canónica de las partes.
// Para /pay, se llama con (saleID, body). Para /sales, solo con (body).
func idempotencyHash(parts ...any) (string, error) {
	h := sha256.New()
	for i, p := range parts {
		b, err := json.Marshal(p)
		if err != nil {
			return "", err
		}
		if i > 0 {
			h.Write([]byte{0}) // separador: evita colisiones tipo ("ab","c") vs ("a","bc")
		}
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}