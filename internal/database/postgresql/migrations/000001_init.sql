-- +goose Up
CREATE TYPE sex_type AS ENUM ('male', 'female');

CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
    first_name VARCHAR(50) NOT NULL,
	second_name VARCHAR(50) NOT NULL,
	birthdate DATE,
	sex sex_type,
	biography VARCHAR(200),
	city VARCHAR(50),
	password VARCHAR(100) NOT NULL
);

-- +goose Down
DROP TABLE users;

DROP TYPE sex_type;