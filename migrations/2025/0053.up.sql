ALTER TABLE tls_certificates ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE tls_certificates ALTER COLUMN created_at DROP DEFAULT;

ALTER TABLE tls_certificates ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE tls_certificates ALTER COLUMN updated_at DROP DEFAULT;
CREATE INDEX index_tls_certificates_on_updated_at ON tls_certificates (updated_at);
