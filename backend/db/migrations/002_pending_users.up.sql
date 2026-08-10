CREATE TABLE pending_users (
	id BIGSERIAL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	username VARCHAR(20) NOT NULL,
	password_hash TEXT NOT NULL,
	token_hash TEXT NOT NULL,
	expires_at TIMESTAMP NOT NULL
);
