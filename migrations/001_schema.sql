-- Tenants (comercios). Frontera de aislamiento multi-tenant.
CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    timezone    TEXT NOT NULL DEFAULT 'UTC'
                CHECK (timezone ~ '^[A-Za-z_]+(/[A-Za-z_+-]+)*$'),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Usuarios. UNIQUE(email) global: un email = un tenant.
CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email          TEXT NOT NULL,
    password_hash  TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (email)
);
CREATE INDEX idx_users_tenant ON users (tenant_id);

-- Catálogo. Dinero en centavos.
CREATE TABLE products (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    price_cents  BIGINT NOT NULL CHECK (price_cents > 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_products_tenant ON products (tenant_id, created_at DESC);

-- Ventas. UNKNOWN es estado de primera clase.
-- sale_date se calcula en la app/SQL con la timezone del tenant (sin DEFAULT).
CREATE TABLE sales (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    status           TEXT NOT NULL DEFAULT 'PENDING'
                     CHECK (status IN ('PENDING','APPROVED','DECLINED','UNKNOWN')),
    total_cents      BIGINT NOT NULL CHECK (total_cents >= 0),
    payment_method   TEXT,
    payment_scenario TEXT CHECK (payment_scenario IN ('APPROVED','DECLINED','TIMEOUT')),
    sale_date        DATE NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_sales_tenant_date_status ON sales (tenant_id, sale_date, status);

-- Líneas. Congelan el precio al momento de la venta.
CREATE TABLE sale_items (
    id                BIGSERIAL PRIMARY KEY,
    sale_id           UUID NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    product_id        UUID NOT NULL REFERENCES products(id),
    quantity          INT NOT NULL CHECK (quantity > 0),
    unit_price_cents  BIGINT NOT NULL CHECK (unit_price_cents >= 0)
);
CREATE INDEX idx_sale_items_sale ON sale_items (sale_id);

-- Intentos de pago. scenario = lo que pidió el cliente; result_status = estado resultante.
CREATE TABLE payment_attempts (
    id               BIGSERIAL PRIMARY KEY,
    sale_id          UUID NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    tenant_id        UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    attempt_number   INT  NOT NULL CHECK (attempt_number > 0),
    scenario         TEXT NOT NULL CHECK (scenario IN ('APPROVED','DECLINED','TIMEOUT')),
    result_status    TEXT NOT NULL CHECK (result_status IN ('APPROVED','DECLINED','UNKNOWN')),
    idempotency_key  TEXT,
    provider_ref     TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (sale_id, attempt_number)
);
CREATE INDEX idx_payment_attempts_sale   ON payment_attempts (sale_id, attempt_number);
CREATE INDEX idx_payment_attempts_tenant ON payment_attempts (tenant_id, created_at DESC);

-- Idempotencia. PK compuesta: misma key en otro endpoint/tenant no colisiona.
CREATE TABLE idempotency_keys (
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    idempotency_key TEXT NOT NULL,
    endpoint        TEXT NOT NULL,
    request_hash    TEXT NOT NULL,
    response_status INT NOT NULL,
    response_body   JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, idempotency_key, endpoint)
);

-- Outbox. BIGSERIAL = versión monotónica para el CAS en Redis.
CREATE TABLE outbox (
    id           BIGSERIAL PRIMARY KEY,
    tenant_id    UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    status       TEXT NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending','processing','published')),
    claimed_at   TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_outbox_pending    ON outbox (id)         WHERE status = 'pending';
CREATE INDEX idx_outbox_processing ON outbox (claimed_at) WHERE status = 'processing';