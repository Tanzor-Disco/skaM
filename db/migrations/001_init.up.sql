CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	username VARCHAR(20),
	email TEXT UNIQUE,
	password_hash TEXT,
	last_year_active SMALLINT
);

CREATE TABLE IF NOT EXISTS rooms (
	id BIGSERIAL PRIMARY KEY,
	room_name VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS room_users (
	id BIGSERIAL PRIMARY KEY,
	room_id BIGINT,
	user_id BIGINT
);

CREATE TABLE IF NOT EXISTS messages (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT,
	room_id BIGINT,
	message_text TEXT
);
