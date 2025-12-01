CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,            -- Reference to the order
    stripe_payment_id VARCHAR(100) NOT NULL, -- Stripe payment ID    
    amount BIGINT NOT NULL,       -- Payment amount
    currency VARCHAR(10) DEFAULT 'IDR',  -- Currency (default: IDR)
    status VARCHAR(20) DEFAULT 'pending', -- Payment status
    created_at TIMESTAMP DEFAULT NOW(),  -- Payment creation time
    updated_at TIMESTAMP DEFAULT NOW(),  -- Last status update
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    UNIQUE (stripe_payment_id)           -- Prevent duplicate Stripe transactions
);