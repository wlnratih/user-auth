/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE IF not exists "user"
(
    id SERIAL PRIMARY KEY,
    name   VARCHAR(60) NOT NULL,
    phone_number VARCHAR(13) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
    );

CREATE OR REPLACE FUNCTION touch()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language plpgsql;

CREATE INDEX IF NOT EXISTS
    user_phone_number ON "user" (phone_number);

CREATE INDEX IF NOT EXISTS
    user_name_phone_number_password ON "user" (name, phone_number, password);

CREATE
OR REPLACE FUNCTION insert_user(phone_number_in VARCHAR(13), name_in VARCHAR(60), password_in VARCHAR(100))
    RETURNS INT
    LANGUAGE SQL as
$$
INSERT INTO "user"(phone_number, name, password)
VALUES (phone_number_in, name_in, password_in)
RETURNING id;
$$;

CREATE
OR REPLACE FUNCTION get_user_by_phone_number(phone_number_in VARCHAR(13))
    RETURNS TABLE(
                     id              INT,
                     name            VARCHAR(60),
                     phone_number    VARCHAR(13),
                     password        VARCHAR(100)
                 )
    LANGUAGE SQL as
$$
SELECT id, name, phone_number, password FROM "user" WHERE phone_number = phone_number_in;
$$;

CREATE TABLE IF not exists "user_login_history"
(
    user_id INT,
    name   VARCHAR(60) NOT NULL,
    phone_number VARCHAR(13) NOT NULL,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE OR REPLACE FUNCTION insert_user_login_history(user_id_in INT, phone_number_in VARCHAR(13), password_in VARCHAR(100), name_in VARCHAR(60))
    RETURNS void
    LANGUAGE SQL as
$$
INSERT INTO user_login_history(user_id, name, phone_number, password)
VALUES (user_id_in, name_in, phone_number_in, password_in);
$$;


CREATE
OR REPLACE FUNCTION get_user_by_id(id_in INT)
    RETURNS TABLE(
                     id              INT,
                     name            VARCHAR(60),
                     phone_number    VARCHAR(13),
                     password        VARCHAR(100)
                 )
    LANGUAGE SQL as
$$
SELECT id, name, phone_number, password FROM "user" WHERE id = id_in;
$$;

CREATE
OR REPLACE FUNCTION update_user(id_in INT, name_in VARCHAR(60), phone_number_in VARCHAR(13))
    RETURNS TABLE(
                     id              INT,
                     name            VARCHAR(60),
                     phone_number    VARCHAR(13),
                     password        VARCHAR(100)
                 )
    LANGUAGE SQL as
$$
UPDATE "user" SET
    name = name_in,
    phone_number = phone_number_in
WHERE id = id_in
RETURNING id, name, phone_number, password;
$$;
