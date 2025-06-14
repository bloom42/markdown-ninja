ALTER TABLE orders ADD COLUMN country TEXT NOT NULL DEFAULT 'XX';
ALTER TABLE orders ALTER COLUMN country DROP DEFAULT;

UPDATE orders SET country = billing_address->>'country_code';
ALTER TABLE orders DROP COLUMN billing_address;

ALTER TABLE orders ADD COLUMN additional_invoice_information TEXT NOT NULL DEFAULT '';
ALTER TABLE orders ALTER COLUMN additional_invoice_information DROP DEFAULT;


ALTER TABLE contacts DROP COLUMN billing_address;
ALTER TABLE contacts RENAME COLUMN country_code TO country;

UPDATE contacts SET country = 'XX' WHERE country = '';
