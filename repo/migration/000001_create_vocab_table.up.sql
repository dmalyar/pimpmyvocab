create table if not exists vocab
(
    id      serial         not null
        constraint vocab_pkey
            primary key,
    user_id integer unique not null
);