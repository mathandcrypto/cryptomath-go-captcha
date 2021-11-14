CREATE TABLE IF NOT EXISTS "captcha_tasks" (
    "uuid" VARCHAR NOT NULL PRIMARY KEY,
    "index" INTEGER NOT NULL,
    "difficulty"  SMALLINT NOT NULL,
    "math" VARCHAR NOT NULL,
    "answer" INTEGER NOT NULL
);

CREATE UNIQUE INDEX index_idx ON "captcha_tasks" USING btree("index");