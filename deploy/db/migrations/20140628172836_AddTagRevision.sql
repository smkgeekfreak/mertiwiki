-- +goose Up
-- SQL in tag_revision 'Up' is executed when this migration is applied
-- +goose StatementBegin
-- Table: tag_revision

CREATE TABLE tag_revision
(
  tag_uid bigint NOT NULL,
  revision_uid bigint NOT NULL,
  status integer,
  created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  modified TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  CONSTRAINT "tag_rev_PK" PRIMARY KEY (tag_uid, revision_uid),
  CONSTRAINT "revision_FK" FOREIGN KEY (revision_uid)
      REFERENCES revision (uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "tag_FK" FOREIGN KEY (tag_uid)
      REFERENCES tag (uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE OR REPLACE VIEW tag_revisiondata AS 
SELECT
	tag_revision.tag_uid, 
	tag_revision.revision_uid, 
	tag_revision.status,
	tag_revision.created,
	tag_revision.modified
FROM
	tag_revision;

-- Function: addTagRevision
CREATE OR REPLACE FUNCTION addTagRevision(_tag_uid bigint, _revision_uid bigint, _status integer)
RETURNS TABLE(
	ret_tag_uid bigint,
	ret_revision_uid bigint,
  	ret_status integer,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 

	BEGIN
		-- 
		RETURN QUERY
		INSERT INTO tag_revision (tag_uid, revision_uid, status) VALUES (_tag_uid, _revision_uid, _status) 
		RETURNING 
			tag_uid, 
			revision_uid, 
			status,
			created,
			modified;
	END;
$func$ language 'plpgsql';

-- Function: deleteTagRevision
CREATE OR REPLACE FUNCTION deleteTagRevision(_tag_uid bigint, _revision_uid bigint, _deletedStatus integer)
RETURNS TABLE(
	ret_tag_uid bigint,
	ret_revision_uid bigint,
  	ret_status integer,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 

	BEGIN
		-- 
		RETURN QUERY
		UPDATE tag_revision tr 
			SET status = _deletedStatus, modified = statement_timestamp() WHERE tag_uid = _tag_uid  AND revision_uid= _revision_uid
		RETURNING 
			tr.tag_uid, 
			tr.revision_uid, 
			tr.status,
			tr.created,
			tr.modified;
	END;
$func$ language 'plpgsql';
--
-- Function: updateTagRevisionStatus
CREATE OR REPLACE FUNCTION updateTagRevisionStatus(_tag_uid bigint, _revision_uid bigint, _status integer)
RETURNS TABLE(
	ret_tag_uid bigint,
	ret_revision_uid bigint,
  	ret_status integer,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 

	BEGIN
		-- 
		RETURN QUERY
		UPDATE tag_revision tr 
			SET status = _status, modified = statement_timestamp() WHERE tr.tag_uid = _tag_uid AND tr.revision_uid = _revision_uid
		RETURNING 
			tr.tag_uid, 
			tr.revision_uid, 
			tr.status,
			tr.created,
			tr.modified;
	END;
$func$ language 'plpgsql';

-- +goose StatementEnd

-- +goose Down
-- SQL tag_revision 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP FUNCTION IF EXISTS addTagRevision(tag_uid bigint, revision_uid bigint, status integer);
DROP FUNCTION IF EXISTS updateTagRevision(tag_uid bigint, revision_uid bigint, status integer);
DROP FUNCTION IF EXISTS deleteTagRevision(tag_uid bigint, revision_uid bigint, deletedStatus integer);
DROP VIEW IF EXISTS tag_revisiondata;
DROP TABLE IF EXISTS tag_revision;
-- +goose StatementEnd


