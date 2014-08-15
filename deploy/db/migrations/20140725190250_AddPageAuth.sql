-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE FUNCTION ownsPage (_auth_uid bigint, _page_uid bigint) 
RETURNS boolean as $func$
     SELECT true 
     from pagedata a_page
     where a_page.uid = $2 and a_page.user_uid = $1;
$func$ language 'sql';

CREATE OR REPLACE FUNCTION authorizeOwnership_page(_auth_uid bigint, _page_uid bigint)
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
		a_page.uid as page_uid,
		a_page.user_uid as author_uid,
		_auth_uid as auth_uid,
		COALESCE(ownsPage(_auth_uid,_page_uid),false) as owns,
		COALESCE(ur.user_rating::bigint,0) as user_rating
	FROM pagedata a_page INNER JOIN revisiondata rd ON rd.page_uid = a_page.uid 
	LEFT OUTER JOIN user_ratingdata ur ON ur.user_uid = _auth_uid
	WHERE a_page.uid = _page_uid; 
	END;
$func$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS ownsPage (bigint, bigint);
DROP FUNCTION IF EXISTS authorizeOwnership_page(bigint, bigint);
-- +goose StatementEnd

