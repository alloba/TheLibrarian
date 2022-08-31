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
    chapter_id    TEXT     not null references chapters,
    relative_path TEXT     not null,
    date_created  datetime not null,
    date_modified datetime not null
);

create table if not exists chapters
(
    id            TEXT     not null primary key unique,
    edition_id    TEXT     not null references editions,
    root_path     TEXT     not null,
    name          TEXT     not null,
    date_created  datetime not null,
    date_modified datetime not null,
    unique (edition_id, name)
);

create table if not exists editions
(
    id             TEXT     not null primary key unique,
    name           TEXT     not null,
    edition_number INT      not null,
    book_id        TEXT     not null references books,
    date_created   datetime not null,
    date_modified  datetime not null
);


