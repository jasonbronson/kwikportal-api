CREATE TABLE users (
    id string PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE TABLE bookmarks (
    id string PRIMARY KEY,
    user_id string,
    folder TEXT,
    url TEXT,
    add_date INTEGER,
    icon TEXT,
    name TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX "bookmark_user_id_url" ON "bookmarks" ("user_id", "url");