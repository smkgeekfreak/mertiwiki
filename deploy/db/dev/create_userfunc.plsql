CREATE OR REPLACE FUNCTION addUser(username varchar, password_hash varchar, email varchar, status integer)
RETURNS void as $$
	INSERT INTO Account (username, password_hash, email, status) VALUES ( $1, $2, $3, $4);
$$ language 'sql';

-- Retrieve a specific user
CREATE OR REPLACE FUNCTION getUser(name varchar)
--RETURNS TABLE(uid bigint, username varchar, password_hash varchar, status int) as $$
--	select uid, username, password_hash, status from Account where username=$1;
RETURNS SETOF account AS $$
SELECT uid, username, password_hash, email, status FROM account;
$$ language 'sql';
