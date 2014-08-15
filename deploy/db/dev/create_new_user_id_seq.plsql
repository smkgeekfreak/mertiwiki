-- Sequence: seq_new_user_id

--DROP SEQUENCE seq_new_user_id;

CREATE SEQUENCE seq_new_user_id
  INCREMENT 1
  MINVALUE 0
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

