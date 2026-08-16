CREATE TABLE pending_users (
	id BIGSERIAL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE CHECK(email <> ''),
	username VARCHAR(20) NOT NULL CHECK(username <> ''),
	password_hash TEXT NOT NULL CHECK(password_hash <> ''),
	token_hash TEXT NOT NULL CHECK(token_hash <> ''),
	expires_at TIMESTAMP NOT NULL
);
