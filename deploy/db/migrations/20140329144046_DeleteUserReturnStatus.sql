
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
DROP FUNCTION IF EXISTS deletePage (pageId integer,deletedStatus integer);
CREATE OR REPLACE FUNCTION deletePage (pageId integer,deletedStatus integer)
RETURNS integer as $func$
	UPDATE Page SET status = $2,  modified = statement_timestamp() WHERE uid = $1 RETURNING status;
$func$ language 'sql'; 
DROP FUNCTION IF EXISTS deleteAccount (accountId integer,deletedStatus integer);
CREATE OR REPLACE FUNCTION deleteAccount (accountId integer, deletedStatus integer)
RETURNS integer as $func$
	UPDATE Account SET status = $2,  modified = statement_timestamp() WHERE uid = $1 RETURNING status;
$func$ language 'sql'; 
-- +goose StatementEnd 

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP FUNCTION IF EXISTS deletePage (pageId integer,deletedStatus integer);
CREATE OR REPLACE FUNCTION deletePage (pageId integer,deletedStatus integer)
RETURNS void as $$
	UPDATE Page SET status = $2,  modified = statement_timestamp() WHERE uid = $1;
$$ language 'sql'; 
DROP FUNCTION IF EXISTS deleteAccount (accountId integer,deletedStatus integer);
CREATE OR REPLACE FUNCTION deleteAccount (accountId integer, deletedStatus integer)
RETURNS void as $$
	UPDATE Account SET status = $2,  modified = statement_timestamp() WHERE uid = $1;
$$ language 'sql'; 
-- +goose StatementEnd 
