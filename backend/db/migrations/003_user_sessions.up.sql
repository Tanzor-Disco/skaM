CREATE TABLE user_sessions (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id),
	session_string TEXT NOT NULL UNIQUE CHECK(session_string <> ''),
	expires_at TIMESTAMP NOT NULL
);
