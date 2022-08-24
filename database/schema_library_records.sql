create table if not exists records
(
    uuid             text    not null,
    hash             text    not null,
    filename         text    not null,
    extension        text    not null,
    version          integer not null,

    literal_location text    not null,
    logical_location text,

    date_entered     text    not null,
    date_created     text    not null,
    date_modified    text
)
