begin;
create table if not exists translation
(
    id             serial  not null
        constraint translation_pkey
            primary key,
    vocab_entry_id integer not null
        constraint translation_vocab_entry_id_fkey
            references vocab_entry,
    text           text    not null,
    class          text    not null,
    position       integer not null
);
create index if not exists translation_vocab_entry_id_index
    on translation (vocab_entry_id);
commit;
