create table if not exists vocab_entry
(
    id            serial      not null
        constraint vocab_entry_pkey
            primary key,
    text          text unique not null,
    transcription text        not null
);
