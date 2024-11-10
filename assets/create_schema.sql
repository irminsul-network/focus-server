create table if not exists point(
    id INTEGER NOT NULL,
    content TEXT NOT NULL,
    encountered INTEGER NOT NULL DEFAULT 1,
    created INTEGER NOT NULL,
    archived INTEGER NOT NULL DEFAULT 0,
    achieved INTEGER NOT NULL DEFAULT 0,
    conquered INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY(id AUTOINCREMENT)
)