-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE FUNCTION ownsRevision (_auth_uid bigint, _revision_uid bigint) 
RETURNS boolean as $func$
     SELECT true 
     from revisiondata rev
     where rev.uid = $2 and rev.user_uid = $1;
$func$ language 'sql';

CREATE OR REPLACE FUNCTION authorizeOwnership_revision(_auth_uid bigint, _revision_uid bigint)
RETURNS TABLE(
	ret_uid bigint,
	ret_author_uid bigint,
	ret_auth_uid bigint,
	ret_owns boolean,
  	ret_user_rating bigint
	)
as $func$ 
	BEGIN
	-- 
	RETURN QUERY
	SELECT DISTINCT
		rev.uid as revision_uid,
		rev.user_uid as author_uid,
		_auth_uid as auth_uid,
		COALESCE(ownsRevision(_auth_uid,_revision_uid),false) as owns,
		COALESCE(ur.user_rating::bigint,0) as user_rating
	FROM revisiondata rev LEFT JOIN user_ratingdata ur ON ur.user_uid = _auth_uid
	WHERE rev.uid = _revision_uid; 	
END;
$func$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS ownsRevision (bigint, bigint);
DROP FUNCTION IF EXISTS authorizeOwnership_revision(bigint, bigint);
-- +goose StatementEnd
