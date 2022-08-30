create table if not exists records
(
    id                 TEXT     not null primary key unique,
    file_pointer       TEXT     not null,
    name               TEXT     not null,
    extension          TEXT     not null,
    date_created       datetime not null,
    date_file_modified datetime not null,
    date_modified      datetime not null
);

create table if not exists books
(
    id            TEXT     not null primary key unique,
    name          TEXT     not null unique,
    date_created  datetime not null,
    date_modified datetime not null
);

create table if not exists pages
(
    id            TEXT     not null primary key unique,
    record_id     TEXT     not null references records,
    edition_id    TEXT     not null references editions,
    date_created  datetime not null,
    date_modified datetime not null
);

create table if not exists editions
(
    id             TEXT     not null primary key unique,
    edition_number INT      not null,
    book_id        TEXT     not null references books,
    date_created   datetime not null,
    date_modified  datetime not null
);


