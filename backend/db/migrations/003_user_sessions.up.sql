CREATE TABLE user_sessions (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id),
	session_id TEXT NOT NULL UNIQUE CHECK(session_id <> ''),
	expires_at TIMESTAMP NOT NULL
);
