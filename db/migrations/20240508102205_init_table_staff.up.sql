CREATE TABLE IF NOT EXISTS staff (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_phone_staff on staff(phone_number);
CREATE INDEX idx_id_staff on staff(id);

CREATE TABLE IF NOT EXISTS country_code (
    id SERIAL PRIMARY KEY,
    phone_number_code VARCHAR(7) UNIQUE NOT NULL,
    country VARCHAR(255) NOT NULL
);

CREATE INDEX idx_phone_number_code_staff on country_code(phone_number_code);