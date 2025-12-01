ALTER TABLE users
DROP COLUMN "password";

ALTER TABLE users
ADD COLUMN password TEXT;