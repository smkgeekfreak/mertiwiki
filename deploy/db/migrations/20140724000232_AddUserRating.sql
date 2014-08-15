-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE OR REPLACE VIEW user_ratingdata AS 
SELECT (rev.user_uid) as user_uid,
	sum(rate.rating)::integer as user_rating,
	statement_timestamp()::TIMESTAMP WITHOUT TIME ZONE as updated 
FROM revision rev LEFT OUTER JOIN rating rate ON rate.revision_uid = rev.uid
WHERE rate.rating IS NOT NULL
GROUP BY rev.user_uid ORDER BY user_rating;

-- Function: findUserRating

CREATE OR REPLACE FUNCTION find_user_rating(_user_uid bigint)
RETURNS TABLE(
	ret_user_uid bigint,
  	ret_user_rating integer,
  	ret_updated TIMESTAMP WITHOUT TIME ZONE
	)
as $func$ 
	BEGIN
		-- 
		RETURN QUERY
		SELECT user_uid, user_rating, updated FROM user_ratingdata
		WHERE user_uid = _user_uid;
	END;
$func$ language 'plpgsql';


-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP FUNCTION IF EXISTS find_user_rating(bigint);
DROP VIEW IF EXISTS user_ratingdata;

-- +goose StatementEnd
