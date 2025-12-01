CREATE TABLE orders (
	id BIGSERIAL PRIMARY KEY,           -- Unique order ID
    user_id BIGINT NOT NULL,             -- User who placed the order
    total_price BIGINT NOT NULL,  -- Total price of the order
    status VARCHAR(20) DEFAULT 'pending', -- Order status
    created_at TIMESTAMP DEFAULT NOW(),  -- Order creation time
    updated_at TIMESTAMP DEFAULT NOW(),  -- Order update time
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);