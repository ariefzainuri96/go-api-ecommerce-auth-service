CREATE TABLE shopping_carts (
	id BIGSERIAL PRIMARY KEY,
	product_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL,
	quantity INT DEFAULT 1 CHECK (quantity > 0),
	created_at TIMESTAMP DEFAULT NOW(),
	FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	UNIQUE (user_id, product_id)
);