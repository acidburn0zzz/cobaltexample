CREATE TABLE jobs (
	id serial PRIMARY KEY,
	title varchar(254) NOT NULL,
	company varchar(254) NOT NULL,
	description text NOT NULL,
	contact varchar(254) NOT NULL,
	created timestamp,
	updated timestamp
);

-- Insert some seed data
INSERT INTO jobs (id, title, company, description, contact, created, updated)
VALUES (1, 'Marketing Developer', 'Acme Co', 'This is a *great* job! Work with our marketing team to develop solutions.', 'Maryam (maryam@example.com)', now(), now()),
(2, 'Designer', 'Acme Co', 'This is _perfect_ for creative types!', 'Maryam (maryam@example.com)', now(), now());
