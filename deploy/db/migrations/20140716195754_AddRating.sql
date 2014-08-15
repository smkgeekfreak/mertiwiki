-- +goose Up
-- SQL in rating 'Up' is executed when this migration is applied
-- +goose StatementBegin
-- Table: rating

CREATE TABLE rating
(
  account_uid bigint NOT NULL,
  revision_uid bigint NOT NULL,
  status integer,
  rating bigint,
  created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  modified TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp(),
  CONSTRAINT "rating_PK" PRIMARY KEY (account_uid, revision_uid),
  CONSTRAINT "revision_FK" FOREIGN KEY (revision_uid)
      REFERENCES revision (uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "account_FK" FOREIGN KEY (account_uid)
      REFERENCES account (uid) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE OR REPLACE VIEW ratingdata AS 
SELECT
	rating.account_uid, 
	rating.revision_uid, 
	rating.status,
	rating.rating,
	rating.created,
	rating.modified
FROM
	rating;

-- Function: addRatingRevision
CREATE OR REPLACE FUNCTION addRating(_account_uid bigint, _revision_uid bigint, _rating bigint)
RETURNS TABLE(
	ret_account_uid bigint,
	ret_revision_uid bigint,
  	ret_rating bigint,
  	ret_created TIMESTAMP WITHOUT TIME ZONE,
	ret_modified TIMESTAMP WITHOUT TIME ZONE 
	)
as $func$ 

	BEGIN
		-- 
		RETURN QUERY
		INSERT INTO rating (account_uid, revision_uid, rating) VALUES (_account_uid, _revision_uid, _rating) 
		RETURNING 
			account_uid, 
			revision_uid, 
			rating,
			created,
			modified;
	END;
$func$ language 'plpgsql';

-- +goose StatementEnd

-- +goose Down
-- SQL rating 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP FUNCTION IF EXISTS addRating(account_uid bigint, revision_uid bigint, rating bigint);
DROP VIEW IF EXISTS ratingdata;
DROP TABLE IF EXISTS rating;
-- +goose StatementEnd


