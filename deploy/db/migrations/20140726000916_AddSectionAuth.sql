-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE FUNCTION ownsSection (_auth_uid bigint, _section_uid bigint) 
RETURNS boolean as $func$
     SELECT true 
     from sectiondata a_section 
     where a_section.uid = $2 and a_section.user_uid = $1;
$func$ language 'sql';

CREATE OR REPLACE FUNCTION authorizeOwnership_section(_auth_uid bigint, _section_uid bigint)
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
		a_section.uid as section_uid,
		a_section.user_uid as author_uid,
		_auth_uid as auth_uid,
		COALESCE(ownsSection(_auth_uid,_section_uid),false) as owns,
		COALESCE(ur.user_rating::bigint,0) as user_rating
	FROM sectiondata a_section LEFT OUTER JOIN user_ratingdata ur ON ur.user_uid = _auth_uid
	WHERE a_section.uid = _section_uid; 
	END;
$func$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS ownsSection (bigint, bigint);
DROP FUNCTION IF EXISTS authorizeOwnership_section(bigint, bigint);
-- +goose StatementEnd

