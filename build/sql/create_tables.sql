CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('employee', 'moderator'))
);

CREATE TABLE pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city VARCHAR(100) NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
    registration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

CREATE TABLE receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    date_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'closed')) DEFAULT 'in_progress',
    created_by UUID REFERENCES users(id)
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id UUID NOT NULL REFERENCES receptions(id),
    type VARCHAR(50) NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    date_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by UUID REFERENCES users(id)
);
