
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
-- Sequence: section_row_id_seq
DROP TABLE IF EXISTS section ;
DROP SEQUENCE IF EXISTS section_row_id_seq;
CREATE SEQUENCE section_row_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

-- Table: section
CREATE TABLE section
(
  uid bigint NOT NULL DEFAULT nextval('section_row_id_seq'::regclass),
  page_uid bigint NOT NULL,
  user_uid bigint NOT NULL,
  name text,
  ordernum integer,
  status integer,
  created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  modified TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  CONSTRAINT "section_PK" PRIMARY KEY (uid),
  CONSTRAINT "page_FK" FOREIGN KEY (page_uid)
      REFERENCES page (uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE OR REPLACE VIEW sectiondata AS 
SELECT
	section.uid, 
	section.page_uid, 
	section.user_uid, 
	section.name,
	section.ordernum,
	section.status,
	section.created,
	section.modified
FROM
	section;

-- Function: create_section
CREATE OR REPLACE FUNCTION create_section (_page_uid bigint, _user_uid bigint, _name text, _ordernum integer, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_name text,
	ret_ordernum integer,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		INSERT INTO section (page_uid, user_uid, name, ordernum, status) 
		VALUES (_page_uid, _user_uid, _name, _ordernum, _status) 
		RETURNING			
			uid, 
			page_uid,
			user_uid,
			name, 
			ordernum,
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';


CREATE OR REPLACE FUNCTION update_section (_uid bigint, _name text, _ordernum integer, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_name text,
	ret_ordernum integer,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE section s 
			SET name = _name, ordernum = _ordernum, status = _status, modified = statement_timestamp() WHERE s.uid = _uid
		RETURNING			
			s.uid, 
			s.page_uid,
			s.user_uid,
			s.name, 
			s.ordernum,
			s.status,
			s.created,
			s.modified;
	END;
$func$ language 'plpgsql';
--Function 
CREATE OR REPLACE FUNCTION update_sectionstatus (_uid bigint, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_name text,
	ret_ordernum integer,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE section s 
			SET status = _status, modified = statement_timestamp() WHERE s.uid = _uid
		RETURNING			
			s.uid, 
			s.page_uid,
			s.user_uid,
			s.name, 
			s.ordernum,
			s.status,
			s.created,
			s.modified;
	END;
$func$ language 'plpgsql';

-- function
CREATE OR REPLACE FUNCTION findbyid_section(_uid bigint)
RETURNS TABLE(
	ret_uid bigint,
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_name text,
	ret_ordernum integer,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 
	BEGIN
		-- 
		RETURN QUERY
		SELECT * FROM sectiondata sd
		WHERE sd.uid = _uid;
	END;
$func$ language 'plpgsql';

CREATE OR REPLACE FUNCTION findbyauthor_section(_user_uid bigint)   
RETURNS SETOF section AS  $BODY$
DECLARE r record;
BEGIN
  FOR r IN SELECT * FROM sectiondata WHERE user_uid = _user_uid 
  LOOP
    RETURN NEXT r;
  END LOOP;
END;
$BODY$ 
LANGUAGE plpgsql;  

-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP FUNCTION IF EXISTS findbyauthor_section(_user_uid bigint);   
DROP FUNCTION IF EXISTS findbyid_section(_uid bigint);
DROP FUNCTION IF EXISTS update_sectionstatus (_uid bigint, _status integer); 
DROP FUNCTION IF EXISTS update_section (_uid bigint, text, _ordernum integer, _status integer);
DROP FUNCTION IF EXISTS create_section (_page_uid bigint, _user_uid bigint, _name text, _ordernum integer, _status integer);

DROP VIEW IF EXISTS sectiondata;
DROP TABLE IF EXISTS section;
DROP SEQUENCE IF EXISTS section_row_id_seq;
-- +goose StatementEnd

