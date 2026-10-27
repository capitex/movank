package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"-"` // interno: nunca se expone
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // interno: nunca se expone
	CreatedAt    time.Time `json:"created_at"`
}

type Product struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"-"` // interno
	Name       string    `json:"name"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
}

type SaleItem struct {
	ProductID      uuid.UUID `json:"product_id"`
	Name           string    `json:"name"` // del JOIN con products (nombre del momento)
	Quantity       int       `json:"quantity"`
	UnitPriceCents int64     `json:"unit_price_cents"` // congelado al momento de la venta
	SubtotalCents  int64     `json:"subtotal_cents"`   // quantity * unit_price_cents
}

type PaymentAttempt struct {
	AttemptNumber int             `json:"attempt_number"`
	Scenario      PaymentScenario `json:"scenario"`
	ResultStatus  PaymentResult   `json:"result_status"`
	ProviderRef   *string         `json:"provider_ref"` // null si TIMEOUT sin respuesta
	CreatedAt     time.Time       `json:"created_at"`
}

type Sale struct {
	ID              uuid.UUID        `json:"id"`
	TenantID        uuid.UUID        `json:"-"`
	Status          SaleStatus       `json:"status"`
	TotalCents      int64            `json:"total_cents"`
	PaymentMethod   string           `json:"payment_method"`
	PaymentScenario *PaymentScenario `json:"payment_scenario"` // null hasta el primer /pay
	SaleDate        Date             `json:"sale_date"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Items           []SaleItem       `json:"items,omitempty"`
	Attempts        []PaymentAttempt `json:"attempts,omitempty"`
}

// NewSaleItem es la entrada para crear una venta. El precio NO se incluye:
// lo resuelve el servidor consultando products por (tenant_id, id).
type NewSaleItem struct {
	ProductID uuid.UUID
	Quantity  int
}

type DashboardAggregate struct {
	Date               Date      `json:"date"`
	ApprovedTotalCents int64     `json:"approved_total_cents"`
	ApprovedCount      int64     `json:"approved_count"`
	UnknownTotalCents  int64     `json:"unknown_total_cents"`
	UnknownCount       int64     `json:"unknown_count"`
	Version            int64     `json:"version"`
	UpdatedAt          time.Time `json:"updated_at"`
}
