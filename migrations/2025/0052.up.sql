ALTER TABLE orders ADD COLUMN country TEXT NOT NULL DEFAULT 'XX';
ALTER TABLE orders ALTER COLUMN country DROP DEFAULT;

UPDATE orders SET country = billing_address->>'country_code';
ALTER TABLE orders DROP COLUMN billing_address;
