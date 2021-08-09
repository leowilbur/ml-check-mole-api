CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE accounts (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	phone VARCHAR(64) NOT NULL,
	gender VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

CREATE TABLE body_parts (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL,
	displayed BOOLEAN NOT NULL,
	image VARCHAR(255) NOT NULL,
	"order" INT NOT NULL
);

CREATE TABLE lesions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	account_id UUID NOT NULL,
	name VARCHAR(255) NOT NULL,
	body_part_id UUID NOT NULL,
	body_part_location VARCHAR(255) NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,

	CONSTRAINT lesions_account_id FOREIGN KEY (account_id)
		REFERENCES accounts (id) ON DELETE CASCADE,
	CONSTRAINT lesions_body_part_id FOREIGN KEY (body_part_id)
		REFERENCES body_parts (id) ON DELETE CASCADE
);

CREATE TABLE questions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL,
	type VARCHAR(64) NOT NULL,
	answers JSONB NOT NULL,
	displayed BOOLEAN NOT NULL,
	"order" INT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

CREATE TABLE requests (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	account_id UUID NOT NULL,
	status VARCHAR(64) NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,

	CONSTRAINT requests_account_id FOREIGN KEY (account_id)
		REFERENCES accounts (id) ON DELETE CASCADE
);

CREATE TABLE reports (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	request_id UUID,
	lesion_id UUID NOT NULL,
	photos VARCHAR (255) ARRAY NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,

	CONSTRAINT reports_request_id FOREIGN KEY (request_id)
		REFERENCES requests (id) ON DELETE CASCADE,
	CONSTRAINT reports_lesion_id FOREIGN KEY (lesion_id)
		REFERENCES lesions (id) ON DELETE CASCADE
);

CREATE TABLE report_answers (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	report_id UUID NOT NULL,
	question_id UUID NOT NULL,
	answer JSONB NOT NULL,

	CONSTRAINT report_answers_report_id FOREIGN KEY (report_id)
		REFERENCES reports (id) ON DELETE CASCADE,
	CONSTRAINT report_answers_question_id FOREIGN KEY (question_id)
		REFERENCES questions (id) ON DELETE CASCADE
);
