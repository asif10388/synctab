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
    username         varchar(256) not null unique check (length(username) >= 1),
    email            varchar(1024) not null unique check (length(email) >= 1),
	passwordHash 	 varchar(1024) not null check (length(passwordHash) >= 1),
    created_at       timestamptz not null default now(),
    updated_at       timestamptz not null default now()
);

/*SQLEND*/

create table if not exists {{.SchemaName}}.urlGroups(
	id 				uuid default gen_random_uuid() primary key,
	user_id 		uuid not null,
	title 			varchar(512) not null check (length(title) >= 1),
	created_at 		timestamptz not null default now(),
	updated_at 		timestamptz not null default now(),
	constraint 		fk_user_id foreign key (user_id) references {{.SchemaName}}.users(id) on delete cascade
)

/*SQLEND*/

create table if not exists {{.SchemaName}}.urls (
	id               uuid default gen_random_uuid() primary key,
	user_id uuid 	 not null,
	url_group_id 	 uuid not null,
	url              varchar(1024) not null check (length(url) >= 1),
	title            varchar(512) not null check (length(title) >= 1),
	created_at       timestamptz not null default now(),
	updated_at       timestamptz not null default now(),
	constraint 		 fk_user_id foreign key (user_id) references {{.SchemaName}}.users(id) on delete cascade,
	constraint 		 fk_url_group_id foreign key (url_group_id) references {{.SchemaName}}.urlGroups(id) on delete cascade
)

/*SQLEND*/