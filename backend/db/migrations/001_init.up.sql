CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	email VARCHAR(300) NOT NULL UNIQUE CHECK(email <> ''),
	username VARCHAR(20) NOT NULL CHECK(username <> ''),
	password_hash TEXT NOT NULL CHECK(password_hash <> ''),
	last_year_active SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
	id BIGSERIAL PRIMARY KEY,
	room_name VARCHAR(20) NOT NULL CHECK(room_name <> '')
);

CREATE TABLE IF NOT EXISTS room_users (
	id BIGSERIAL PRIMARY KEY,
	room_id BIGINT NOT NULL,
	user_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL,
	room_id BIGINT NOT NULL,
	message_text TEXT NOT NULL CHECK(message_text <> '')
);
