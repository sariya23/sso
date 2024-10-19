create table if not exists "user" (
    user_id bigint generated always as identity primary key,
    email varchar(320) not null unique,
    pass_hash bytea not null
);

create table if not exists app (
    app_id smallint primary key,
    name varchar(100) not null unique,
    secret text not null unique
);