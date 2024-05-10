CREATE TABLE IF NOT EXISTS order_detail (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    product_order_quantity INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    -- ,updated_at TIMESTAMP
);

CREATE INDEX idx_id_order_detail on order_detail(id);
CREATE INDEX idx_order_id_order_detail on order_detail(order_id);
CREATE INDEX idx_product_id_order_detail on order_detail(product_id);
CREATE INDEX idx_created_at_desc_order_detail ON order_detail(created_at DESC);