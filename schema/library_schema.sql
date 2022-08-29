create table if not exists record
(
    hash               TEXT     not null primary key unique,
    file_pointer       TEXT     not null,
    name               TEXT     not null,
    extension          TEXT     not null,
    date_created       datetime not null,
    date_file_modified datetime not null,
    date_modified      datetime not null
);

create table if not exists book
(
    uuid          TEXT     not null primary key unique,
    name          TEXT     not null unique,
    date_created  datetime not null,
    date_modified datetime not null
);

create table if not exists page
(
    uuid          TEXT     not null primary key unique,
    record_hash   TEXT     not null references record,
    book_uuid     TEXT     not null references book,
    edition_uuid  TEXT     not null references edition,
    date_created  datetime not null,
    date_modified datetime not null
);

create table if not exists edition
(
    uuid           TEXT     not null primary key unique,
    edition_number INT      not null,
    book_uuid      TEXT     not null references book,
    date_created   datetime not null,
    date_modified  datetime not null
);


