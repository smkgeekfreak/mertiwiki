-- +goose Up
-- SQL in tag 'Up' is executed when this migration is applied
-- +goose StatementBegin
-- Sequence: tag_row_id_seq
DROP TABLE IF EXISTS tag ;
DROP SEQUENCE IF EXISTS tag_row_id_seq;
CREATE SEQUENCE tag_row_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

-- Table: tag
CREATE TABLE tag
(
  uid bigint NOT NULL DEFAULT nextval('tag_row_id_seq'::regclass),
  name text,
  description text,
  status integer,
  created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  modified TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  CONSTRAINT "tag_PK" PRIMARY KEY (uid)
);

CREATE OR REPLACE VIEW tagdata AS 
SELECT
	tag.uid, 
	tag.name, 
	tag.description, 
	tag.status,
	tag.created,
	tag.modified
FROM
	tag;

-- Function: create_tag
CREATE OR REPLACE FUNCTION create_tag(_name text, _description text, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_name text,
  	ret_description text,
  	ret_status integer,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 

	BEGIN
		-- 
		RETURN QUERY
			INSERT INTO tag (name, description, status) VALUES (_name, _description, _status)
		RETURNING 
			uid, 
			name, 
			description, 
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';

-- Function: updateTag
CREATE OR REPLACE FUNCTION update_tag(_uid bigint, _name text, _description text, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_name text,
  	ret_description text,
  	ret_status integer,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 

	BEGIN
		-- 
		RETURN QUERY
		UPDATE tag t 
			SET name = _name, description = _description, status = _status, modified = statement_timestamp() WHERE t.uid = _uid 
		RETURNING 
			t.uid, 
			t.name, 
			t.description, 
			t.status,
			t.created,
			t.modified;
	END;
$func$ language 'plpgsql';

-- Function: updateTagStatus

CREATE OR REPLACE FUNCTION update_tagstatus(_uid bigint, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_name text,
  	ret_description text,
  	ret_status integer,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 

	BEGIN
		-- 
		RETURN QUERY
		UPDATE tag t 
			SET status = _status, modified = statement_timestamp() WHERE t.uid = _uid
		RETURNING 
			t.uid, 
			t.name, 
			t.description, 
			t.status,
			t.created,
			t.modified;
	END;
$func$ language 'plpgsql';

-- function
CREATE OR REPLACE FUNCTION findbyid_tag(_uid bigint)
RETURNS TABLE(
	ret_uid bigint,
	ret_name text,
  	ret_description text,
  	ret_status integer,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 
	BEGIN
		-- 
		RETURN QUERY
		SELECT * FROM tagdata td
		WHERE td.uid = _uid;
	END;
$func$ language 'plpgsql';

-- +goose StatementEnd

-- +goose Down
-- SQL tag 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP FUNCTION IF EXISTS create_tag(name text, description text, status integer);
DROP FUNCTION IF EXISTS update_tag(uid bigint, name text, description text, status integer);
DROP FUNCTION IF EXISTS update_tagstatus(uid bigint, status integer);
DROP VIEW IF EXISTS tagdata;
DROP TABLE IF EXISTS tag;
DROP SEQUENCE IF EXISTS tag_row_id_seq;
-- +goose StatementEnd



