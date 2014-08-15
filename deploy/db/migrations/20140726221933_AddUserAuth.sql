-- +goose Up
-- SQL in user 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE FUNCTION isAccount(_auth_uid bigint, _user_uid bigint) 
RETURNS boolean as $func$
     SELECT true 
     from userdata a_user 
     where a_user.uid = $2 and $2 = $1;
$func$ language 'sql';

CREATE OR REPLACE FUNCTION authorizeOwnership_account(_auth_uid bigint, _user_uid bigint)
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
		a_user.uid as user_uid,
		a_user.uid as author_uid,
		_auth_uid as auth_uid,
		COALESCE(isAccount(_auth_uid,_user_uid),false) as owns,
		COALESCE(ur.user_rating::bigint,0) as user_rating
	FROM userdata a_user LEFT OUTER JOIN user_ratingdata ur ON ur.user_uid = _auth_uid
	WHERE a_user.uid = _user_uid; 
	END;
$func$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS isAccount(bigint, bigint);
DROP FUNCTION IF EXISTS authorizeOwnership_account(bigint, bigint);
-- +goose StatementEnd

