CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email on users(lower(email))
