CREATE TABLE IF NOT EXISTS urls (
    short_code VARCHAR(10) PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(original_url)
);

CREATE INDEX IF NOT EXISTS idx_urls_original_url ON urls(original_url);
CREATE INDEX IF NOT EXISTS idx_urls_expires_at ON urls(expires_at);

DROP TABLE IF EXISTS urls; 