create table if not exists point(
    id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_on INTEGER NOT NULL,
    archived INTEGER NOT NULL DEFAULT 0,
    achieved INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY(id AUTOINCREMENT)
);

create table if not exists point_encounters (
    id integer not null PRIMARY KEY,
    point_id integer not null,
    encountered_on integer not null,
    conquered real not null default 0,
    urgency integer not null default 1,
    foreign key(point_id) references point(id)
);