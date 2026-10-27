package domain

// ── SaleStatus ────────────────────────────────────────────────────────

type SaleStatus string

const (
	SaleStatusPending  SaleStatus = "PENDING"
	SaleStatusApproved SaleStatus = "APPROVED"
	SaleStatusDeclined SaleStatus = "DECLINED"
	SaleStatusUnknown  SaleStatus = "UNKNOWN"
)

func (s SaleStatus) IsValid() bool {
	switch s {
	case SaleStatusPending, SaleStatusApproved, SaleStatusDeclined, SaleStatusUnknown:
		return true
	}
	return false
}

// IsTerminal indica si la venta ya no puede moverse por /pay ni por el resolver.
// APPROVED y DECLINED son terminales. UNKNOWN no: lo mueve el resolver.
func (s SaleStatus) IsTerminal() bool {
	return s == SaleStatusApproved || s == SaleStatusDeclined
}

// IsPayable indica si /pay puede actuar sobre este estado.
// Solo PENDING. UNKNOWN se resuelve por el resolver interno, nunca por /pay.
func (s SaleStatus) IsPayable() bool {
	return s == SaleStatusPending
}

// AffectsDashboardBetween indica si la transición old → new cambia algún
// cubo del dashboard. Regla del punto 4: se emite outbox si y solo si
// el estado entra, sale o se mueve entre {APPROVED, UNKNOWN}.
func AffectsDashboardBetween(old, new SaleStatus) bool {
	return dashboardBucket(old) != dashboardBucket(new)
}

func dashboardBucket(s SaleStatus) string {
	switch s {
	case SaleStatusApproved:
		return "approved"
	case SaleStatusUnknown:
		return "unknown"
	default:
		return ""
	}
}

// ── PaymentScenario ───────────────────────────────────────────────────

type PaymentScenario string

const (
	PaymentScenarioApproved PaymentScenario = "APPROVED"
	PaymentScenarioDeclined PaymentScenario = "DECLINED"
	PaymentScenarioTimeout  PaymentScenario = "TIMEOUT"
)

func (p PaymentScenario) IsValid() bool {
	switch p {
	case PaymentScenarioApproved, PaymentScenarioDeclined, PaymentScenarioTimeout:
		return true
	}
	return false
}

// ResolveStatus mapea el scenario pedido por el cliente al resultado del pago.
// Regla dura de la prueba: TIMEOUT siempre produce UNKNOWN, nunca DECLINED.
// Centralizar la regla aquí garantiza que no hay otra ruta que la mapee mal.
func (p PaymentScenario) ResolveStatus() (PaymentResult, error) {
	switch p {
	case PaymentScenarioApproved:
		return PaymentResultApproved, nil
	case PaymentScenarioDeclined:
		return PaymentResultDeclined, nil
	case PaymentScenarioTimeout:
		return PaymentResultUnknown, nil
	default:
		return "", ErrInvalidScenario
	}
}

// ── PaymentResult ─────────────────────────────────────────────────────

type PaymentResult string

const (
	PaymentResultApproved PaymentResult = "APPROVED"
	PaymentResultDeclined PaymentResult = "DECLINED"
	PaymentResultUnknown  PaymentResult = "UNKNOWN"
)

// ToSaleStatus convierte el resultado del pago al estado resultante de la venta.
// Los strings coinciden por diseño; el cast es seguro.
func (r PaymentResult) ToSaleStatus() SaleStatus {
	return SaleStatus(r)
}
