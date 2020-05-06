begin;
create table if not exists vocab_to_entry_link
(
    vocab_id integer not null
        constraint vocab_to_entry_link_vocab_id_fkey
            references vocab,
    entry_id integer not null
        constraint vocab_to_entry_link_entry_id_fkey
            references vocab_entry
);
create unique index if not exists vocab_to_entry_link_ids_uindex
    on vocab_to_entry_link (vocab_id, entry_id);
commit;
