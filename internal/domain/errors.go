package domain

import "errors"

var (
	// ErrNotFound: recurso inexistente o perteneciente a otro tenant.
	// El handler HTTP lo mapea a 404 (nunca 403: no filtramos existencia).
	ErrNotFound = errors.New("not found")

	// ErrSaleNotPending: /pay sobre una venta que no está PENDING.
	// Mapea a 409 SALE_NOT_PENDING.
	ErrSaleNotPending = errors.New("sale not pending")

	// ErrProductNotFound: product_id inexistente o de otro tenant.
	// Mapea a 422 PRODUCT_NOT_FOUND.
	ErrProductNotFound = errors.New("product not found")

	// ErrInvalidScenario: scenario fuera del enum permitido.
	// Mapea a 422 INVALID_SCENARIO.
	ErrInvalidScenario = errors.New("invalid scenario")

	// ErrIdempotencyConflict: misma key + mismo endpoint + body distinto.
	// Mapea a 409 IDEMPOTENCY_CONFLICT.
	ErrIdempotencyConflict = errors.New("idempotency conflict")
)

// ValidationError es un error de validación del body con el campo que falló.
// Mapea a 400 INVALID_BODY. Se usa errors.As para detectarlo.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field == "" {
		return "invalid input: " + e.Message
	}
	return "invalid input: " + e.Field + ": " + e.Message
}

// NewValidationError construye un ValidationError como error.
func NewValidationError(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}
