CREATE TABLE IF NOT EXISTS order_product (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    paid INT NOT null,
    change INT NOT null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    -- ,updated_at TIMESTAMP
);

CREATE INDEX idx_id_order on order_product(id);
CREATE INDEX idx_customer_id_order_product on order_product(customer_id);
CREATE INDEX idx_created_at_desc_order_product ON order_product(created_at DESC);