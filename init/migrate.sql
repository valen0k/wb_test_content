create table weather
(
    city varchar primary key not null,
    temp double precision    not null,
    date timestamp default current_timestamp,
    data jsonb               not null
);
