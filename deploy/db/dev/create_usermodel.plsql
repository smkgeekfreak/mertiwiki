-- Table: "User_test"

CREATE TABLE Account 
(
  uid bigint NOT NULL DEFAULT nextval('seq_new_user_id'::regclass),
  username character varying(50),
  password_hash character varying(255),
  email character varying(255),
  status integer,
  CONSTRAINT PK_USER_TEST PRIMARY KEY (uid)
);


--
-- Create View of UserData from Account
--
CREATE OR REPLACE VIEW userdata AS 
 SELECT account.uid, 
    account.username, 
    account.password_hash, 
    account.email, 
    account.status
   FROM account;


