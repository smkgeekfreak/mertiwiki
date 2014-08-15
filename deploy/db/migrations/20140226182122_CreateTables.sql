
-- +goose Up
-- +goose StatementBegin
DROP SEQUENCE IF EXISTS account_row_id_seq;
CREATE SEQUENCE account_row_id_seq CYCLE;

-- place unique constraint on account username
CREATE TABLE Account 
(
  uid bigint NOT NULL DEFAULT nextval('account_row_id_seq'::regclass) PRIMARY KEY,
  username character varying(50),
  password_hash character varying(255),
  email character varying(255),
  status integer,
  created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  modified TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()
);

CREATE OR REPLACE VIEW userdata AS 
SELECT
	Account.uid, 
	Account.username, 
	Account.password_hash, 
	Account.email, 
	Account.status,
	Account.created,
	Account.modified
FROM
	Account;

-- Function: create_account
CREATE OR REPLACE FUNCTION create_account(_username varchar, _password_hash varchar, _email varchar, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_username varchar,
	ret_password_hash varchar,
	ret_email varchar,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		INSERT INTO account(username, password_hash, email, status) 
		VALUES (_username, _password_hash, _email, _status) 
		RETURNING			
			uid, 
			username, 
			password_hash, 
			email, 
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_account(_uid bigint, _email varchar, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_username varchar,
	ret_password_hash varchar,
	ret_email varchar,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE account a 
			SET email = _email, status = _status, modified = statement_timestamp() WHERE a.uid = _uid
		RETURNING			
			uid, 
			username, 
			password_hash, 
			email, 
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_accountstatus(_uid bigint, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_username varchar,
	ret_password_hash varchar,
	ret_email varchar,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE account a 
			SET status = _status, modified = statement_timestamp() WHERE a.uid = _uid
		RETURNING			
			uid, 
			username, 
			password_hash, 
			email, 
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';

CREATE OR REPLACE FUNCTION findbyid_account(_uid bigint)
RETURNS TABLE(
	ret_uid bigint,
	ret_username varchar,
	ret_password_hash varchar,
	ret_email varchar,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 
	BEGIN
		-- 
		RETURN QUERY
		SELECT * FROM userdata ud
		WHERE ud.uid = _uid;
	END;
$func$ language 'plpgsql';

-- PAGE MODEL
DROP SEQUENCE IF EXISTS page_row_id_seq;
CREATE SEQUENCE page_row_id_seq CYCLE;

CREATE TABLE page
(
  uid bigint NOT NULL DEFAULT nextval('page_row_id_seq'::regclass),
  user_uid bigint NOT NULL,
  title text,
  status integer,
  created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  modified TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  CONSTRAINT "page_PK" PRIMARY KEY (uid),
  CONSTRAINT "user_FK" FOREIGN KEY (user_uid)
      REFERENCES account(uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE OR REPLACE VIEW pagedata AS 
SELECT
	Page.uid, 
	Page.user_uid,
	Page.title, 
	Page.status,
	Page.created,
	Page.modified
FROM
	Page;

-- Function: create_page
CREATE OR REPLACE FUNCTION create_page (_user_uid bigint, _title text, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_user_uid bigint,
	ret_title text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		INSERT INTO page (user_uid, title, status) 
		VALUES (_user_uid, _title, _status) 
		RETURNING			
			uid, 
			user_uid,
			title, 
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';


CREATE OR REPLACE FUNCTION update_page (_uid bigint, _title text, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_user_uid bigint,
	ret_title text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE page p 
			SET title = _title, status = _status, modified = statement_timestamp() WHERE p.uid = _uid
		RETURNING			
			p.uid, 
			p.user_uid,
			p.title, 
			p.status,
			p.created,
			p.modified;
	END;
$func$ language 'plpgsql';
--Function 
CREATE OR REPLACE FUNCTION update_pagestatus (_uid bigint, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_user_uid bigint,
	ret_title text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE page p 
			SET status = _status, modified = statement_timestamp() WHERE p.uid = _uid
		RETURNING			
			p.uid, 
			p.user_uid,
			p.title, 
			p.status,
			p.created,
			p.modified;
	END;
$func$ language 'plpgsql';

-- function
CREATE OR REPLACE FUNCTION findbyid_page(_uid bigint)
RETURNS TABLE(
	ret_uid bigint,
	ret_user_uid bigint,
	ret_title text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 
	BEGIN
		-- 
		RETURN QUERY
		SELECT * FROM pagedata pd
		WHERE pd.uid = _uid;
	END;
$func$ language 'plpgsql';

CREATE OR REPLACE FUNCTION findbyauthor_page(_user_uid bigint)   
RETURNS SETOF page AS  $BODY$
DECLARE r record;
BEGIN
  FOR r IN SELECT * FROM pagedata WHERE user_uid = _user_uid
  LOOP
    RETURN NEXT r;
  END LOOP;
END;
$BODY$ 
LANGUAGE plpgsql;  

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP FUNCTION IF EXISTS findbyauthor_page(_user_uid bigint);   
DROP FUNCTION IF EXISTS findbyid_page(_uid bigint);
DROP FUNCTION IF EXISTS update_pagestatus (_uid bigint, _status integer); 
DROP FUNCTION IF EXISTS update_page (_uid bigint, _title text, _status integer);
DROP FUNCTION IF EXISTS create_page (_user_uid bigint, _title text, _status integer);
DROP VIEW IF EXISTS pagedata;
DROP TABLE IF EXISTS page;
DROP SEQUENCE IF EXISTS page_row_id_seq;

DROP FUNCTION IF EXISTS create_account(_username varchar, _password_hash varchar, _email varchar, _status integer);
DROP FUNCTION IF EXISTS update_account(_uid bigint, _email varchar, _status integer);
DROP FUNCTION IF EXISTS update_accountstatus(_uid bigint,  _status integer);
DROP FUNCTION IF EXISTS findbyid_account(_uid bigint);
DROP VIEW IF EXISTS userdata;
DROP TABLE IF EXISTS Account;
DROP SEQUENCE IF EXISTS account_row_id_seq;
-- +goose StatementEnd
