CREATE TABLE IF NOT EXISTS bookmarks
(
    id varchar(36) unique,
    description varchar(255),
    url varchar(2048) not null,
    code varchar(10) not null,
    user_id varchar(36) not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT pkey_bookmarks PRIMARY KEY (id),
    CONSTRAINT unique_code UNIQUE (code),
    CONSTRAINT fk_users_bookmarks FOREIGN KEY (user_id) REFERENCES users(id)
);