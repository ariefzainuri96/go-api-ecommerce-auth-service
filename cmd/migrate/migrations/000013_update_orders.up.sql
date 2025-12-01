ALTER TABLE orders
ADD COLUMN invoice_id BIGINT UNIQUE;

ALTER TABLE orders
ADD COLUMN invoice_url TEXT;

ALTER TABLE orders
ADD COLUMN invoice_exp_date TEXT;