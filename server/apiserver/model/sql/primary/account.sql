DO
$$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = 'main') THEN
        CREATE SCHEMA main;
    END IF;
END
$$;

/*SQLEND*/

create table if not exists {{.SchemaName}}.users (
    id               uuid default gen_random_uuid() primary key,
    created_at       timestamptz not null default now(),
    updated_at       timestamptz not null default now(),
    email            varchar(1024) not null unique check (length(email) >= 1),
    user_id         varchar(256) not null unique check (length(login_id) >= 1),
    active           boolean not null,
    details          jsonb,
);