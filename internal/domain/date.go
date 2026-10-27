package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

// Date representa una fecha sin hora (columna DATE de Postgres).
// JSON: "2006-01-02". Scan: acepta time.Time (pgx / database/sql).
type Date struct {
	time.Time
}

func NewDate(t time.Time) Date {
	y, m, d := t.Date()
	return Date{time.Date(y, m, d, 0, 0, 0, 0, time.UTC)}
}

func (d Date) String() string { return d.Format(dateLayout) }

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("Date.UnmarshalJSON: %w", err)
	}
	d.Time = t
	return nil
}

// Scan implementa sql.Scanner (pgx v5 lo reconoce vía compat layer).
func (d *Date) Scan(src any) error {
	t, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("Date.Scan: esperaba time.Time, recibí %T", src)
	}
	d.Time = t
	return nil
}

// Value implementa driver.Valuer para escribir la fecha como DATE.
func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}
