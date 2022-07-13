CREATE TABLE users (
	id serial NOT NULL,
	auth0_id text NULL,
	first_name text NULL,
	last_name text NULL,
	email text NULL,
	"password" text NULL,
	phone_number text NULL,
	gender int4 NULL,
	username text NULL,
	date_of_birth numeric NULL,
	biography text NULL,
	education text NULL,
	work_experience text NULL,
	skills text NULL,
	interests text NULL,
	active bool NULL,
	public bool NULL,
	message_notifications bool NULL,
	follow_notifications bool NULL,
	like_notifications bool NULL,
	comment_notifications bool NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);

INSERT INTO users
(id, auth0_id, first_name, last_name, email, "password", phone_number, gender, username, date_of_birth, biography, education, work_experience, skills, interests, active, public, message_notifications, follow_notifications, like_notifications, comment_notifications)
VALUES(1, 'auth0|62af383e504e5680df88c742', 'Petar', 'Petrovic', 'admin@dislinkt.com', '$2a$10$GNysTh1mfPQbnNUHQM.iCe5cLIejAWU.6A1TTPDUOa/3.aUvlyG3a', '060123456', 0, 'admin', 315529199616, '', '', '', '', '', false, false, false, false, false, false);
INSERT INTO users
(id, auth0_id, first_name, last_name, email, "password", phone_number, gender, username, date_of_birth, biography, education, work_experience, skills, interests, active, public, message_notifications, follow_notifications, like_notifications, comment_notifications)
VALUES(2, 'auth0|62af385cb690199c1c89faab', 'Laza', 'Lazic', 'admin2@dislinkt.com', '$2a$10$GNysTh1mfPQbnNUHQM.iCe5cLIejAWU.6A1TTPDUOa/3.aUvlyG3a', '060123457', 0, 'admin2', 315529199616, '', '', '', '', '', false, false, false, false, false, false);
INSERT INTO users
(id, auth0_id, first_name, last_name, email, "password", phone_number, gender, username, date_of_birth, biography, education, work_experience, skills, interests, active, public, message_notifications, follow_notifications, like_notifications, comment_notifications)
VALUES(3, 'auth0|62af387270e7f4c2c978fbc4', 'Mita', 'Mitic', 'admin3@dislinkt.com', '$2a$10$GNysTh1mfPQbnNUHQM.iCe5cLIejAWU.6A1TTPDUOa/3.aUvlyG3a', '060123458', 0, 'admin3', 315529199616, '', '', '', '', '', false, false, false, false, false, false);
