CREATE TYPE category_product_name AS ENUM ('Clothing', 'Accessories', 'Footwear', 'Beverages');

CREATE TABLE IF NOT EXISTS product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    sku VARCHAR(30) NOT NULL,
    category category_product_name NOT NULL,
    imageUrl VARCHAR(255) NOT null,
    notes VARCHAR(200) NOT null,
    price INT NOT null,
    stock INT NOT null CHECK (stock BETWEEN 0 AND 100000),
    location VARCHAR(200) NOT null,
    is_available BOOLEAN NOT null,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_id_product on product(id);
CREATE UNIQUE INDEX idx_sku_product_unique on product(sku);
CREATE INDEX idx_created_at_desc_deleted_at_null_product ON product (created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_stock_created_at_desc_deleted_at_null_product ON product (stock ASC, created_at DESC) WHERE deleted_at IS NULL AND is_available = True;