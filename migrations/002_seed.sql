-- Source - https://stackoverflow.com/a/71461372
-- Posted by nbk, modified by community. See post 'Timeline' for change history
-- Retrieved 2026-10-27, License - CC BY-SA 4.0

CREATE EXTENSION pgcrypto;

-- Dos tenants para aislamiento multi-tenant.
INSERT INTO tenants (id, name, timezone)
VALUES (
        '00000000-0000-0000-0000-000000000001',
        'Tenant Demo 1',
        'America/Bogota'
    ),
    (
        '00000000-0000-0000-0000-000000000002',
        'Tenant Demo 1',
        'America/Mexico_City'
    ) ON CONFLICT (id) DO NOTHING;

-- Contraseña de ambos usuarios demo: demo
INSERT INTO users (id, tenant_id, email, password_hash)
VALUES (
        '10000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000001',
        'usuario1@email.test',
        crypt('demo', gen_salt('bf',10))
    ),
    (
        '10000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000002',
        'usuario2@email.test',
        crypt('demo', gen_salt('bf',10))
    ) ON CONFLICT (email) DO NOTHING;

INSERT INTO products (id, tenant_id, name, price_cents)
VALUES (
        '20000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000001',
        'P1',
        350000
    ),
    (
        '20000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000001',
        'P2',
        1200000
    ),
    (
        '20000000-0000-0000-0000-000000000003',
        '00000000-0000-0000-0000-000000000001',
        'P3',
        500000
    ),
    (
        '20000000-0000-0000-0000-000000000004',
        '00000000-0000-0000-0000-000000000001',
        'P4',
        400000
    ),
    (
        '20000000-0000-0000-0000-000000000005',
        '00000000-0000-0000-0000-000000000002',
        'P5',
        800000
    ),
    (
        '20000000-0000-0000-0000-000000000006',
        '00000000-0000-0000-0000-000000000002',
        'P6',
        900000
    ) ON CONFLICT (id) DO NOTHING;