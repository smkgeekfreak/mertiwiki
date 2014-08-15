-- +goose Up
-- SQL in revision 'Up' is executed when this migration is applied
-- +goose StatementBegin
-- Sequence: revision_row_id_seq
DROP TABLE IF EXISTS revision ;
DROP SEQUENCE IF EXISTS revision_row_id_seq;
CREATE SEQUENCE revision_row_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

-- Table: revision
CREATE TABLE revision
(
  uid bigint NOT NULL DEFAULT nextval('revision_row_id_seq'::regclass),
  sec_uid bigint NOT NULL,
  page_uid bigint NOT NULL,
  user_uid bigint NOT NULL,
  body text,
  status integer,
  created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  modified TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  CONSTRAINT "revision_PK" PRIMARY KEY (uid),
  CONSTRAINT "sec_FK" FOREIGN KEY (sec_uid )
      REFERENCES section (uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "user_FK" FOREIGN KEY (user_uid)
      REFERENCES account(uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE OR REPLACE VIEW revisiondata AS 
SELECT
	revision.uid, 
	revision.sec_uid, 
	revision.page_uid, 
	revision.user_uid, 
	revision.body,
	revision.status,
	revision.created,
	revision.modified
FROM
	revision;

-- Function: addRevision
CREATE OR REPLACE FUNCTION create_Revision(_sec_uid bigint, _user_uid bigint, _body text, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_sec_uid bigint, 
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_body text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		INSERT INTO revision (sec_uid, page_uid, user_uid, body, status) 
		VALUES (_sec_uid, (SELECT page_uid FROM sectiondata sd WHERE sd.uid = _sec_uid), _user_uid, _body, _status) 
		RETURNING			
			uid, 
			sec_uid, 
			page_uid, 
			user_uid,
			body, 
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';

-- Function: updateRevision
CREATE OR REPLACE FUNCTION update_Revision(_uid bigint, _body text, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_sec_uid bigint, 
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_body text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE revision r 
			SET body= _body, status = _status, modified = statement_timestamp() WHERE r.uid = _uid
		RETURNING			
			r.uid, 
			r.sec_uid, 
			r.page_uid, 
			r.user_uid,
			r.body, 
			r.status,
			r.created,
			r.modified;
	END;
$func$ language 'plpgsql';

-- Function: updateRevision
CREATE OR REPLACE FUNCTION update_RevisionStatus(_uid bigint, _status integer)
RETURNS TABLE(
	ret_uid bigint,
	ret_sec_uid bigint, 
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_body text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 

	BEGIN
		-- 
	RETURN QUERY
		UPDATE revision r 
			SET status = _status, modified = statement_timestamp() WHERE r.uid = _uid
		RETURNING			
			r.uid, 
			r.sec_uid, 
			r.page_uid, 
			r.user_uid,
			r.body, 
			r.status,
			r.created,
			r.modified;
	END;
$func$ language 'plpgsql';

CREATE OR REPLACE FUNCTION findbyid_revision(_revision_uid bigint)
RETURNS TABLE(
	ret_uid bigint,
	ret_sec_uid bigint, 
	ret_page_uid bigint,
	ret_user_uid bigint,
	ret_body text,
	ret_status integer,
	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 
	BEGIN
		-- 
		RETURN QUERY
		SELECT * FROM revisiondata rd
		WHERE rd.uid = _revision_uid;
	END;
$func$ language 'plpgsql';

CREATE OR REPLACE FUNCTION findbyauthor_revision(userId bigint)   
RETURNS SETOF revision AS  $BODY$
DECLARE r record;
BEGIN
  FOR r IN SELECT * FROM revisiondata WHERE user_uid = userId 
  LOOP
    RETURN NEXT r;
  END LOOP;
END;
$BODY$ 
LANGUAGE plpgsql;  

-- +goose StatementEnd

-- +goose Down
-- SQL revision 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP FUNCTION IF EXISTS findbyauthor_revision(userId bigint);
DROP FUNCTION IF EXISTS findbyid_revision(_revision_uid bigint);
DROP FUNCTION IF EXISTS create_Revision(sec_uid bigint, user_uid bigint, body text, status integer);
DROP FUNCTION IF EXISTS update_Revision(uid bigint, body text, status integer);
DROP FUNCTION IF EXISTS update_RevisionStatus(uid integer, status integer);
DROP VIEW IF EXISTS revisiondata;
DROP TABLE IF EXISTS revision;
DROP SEQUENCE IF EXISTS revision_row_id_seq;
-- +goose StatementEnd


