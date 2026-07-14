CREATE TABLE users
(
    id              uuid unique,
    display_name    varchar(255) not null,
    username        varchar(255) not null,
    email           varchar(255) not null,
    password        varchar(255) not null,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE ,
    CONSTRAINT pkey_users PRIMARY KEY (id),
    CONSTRAINT unique_username UNIQUE (username),
    CONSTRAINT unique_email UNIQUE (email)
);
