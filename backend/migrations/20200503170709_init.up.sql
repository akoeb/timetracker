CREATE TABLE projects(
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "name" VARCHAR NOT NULL,
        "client_name" VARCHAR,
        "status" VARCHAR NOT NULL
    );

CREATE TABLE events (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "project_id" INTEGER NOT NULL REFERENCES projects(id),
        "code" VARCHAR NOT NULL,
        "timestamp" TIMESTAMP NOT NULL,
        "note" TEXT
);