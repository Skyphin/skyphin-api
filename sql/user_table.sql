CREATE TABLE users (
    id bigint primary key generated always as identity,
    username text NOT NULL,
    email text NOT NULL,
    encrypted_password text,
    verified boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT unique_email UNIQUE (email),
    CONSTRAINT unique_username UNIQUE (username)
) WITH (OIDS=FALSE);