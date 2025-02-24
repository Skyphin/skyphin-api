CREATE TABLE reset_tokens (
    id bigint primary key generated always as identity,
    user_id bigint NOT NULL,
    token text NOT NULL,
    expires_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) WITH (OIDS=FALSE);