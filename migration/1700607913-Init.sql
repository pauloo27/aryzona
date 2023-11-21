-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- user table
CREATE TABLE public."user" (
	id varchar(255) NOT NULL,
	preferred_locale varchar(5) NULL,
	last_slash_command_locale varchar(5) NULL,
	CONSTRAINT user_pkey PRIMARY KEY (id)
);

-- guild table
CREATE TABLE public.guild (
	id varchar(255) NOT NULL,
	preferred_locale varchar(5) NOT NULL,
	CONSTRAINT guild_pkey PRIMARY KEY (id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE public."user";
DROP TABLE public.guild;
