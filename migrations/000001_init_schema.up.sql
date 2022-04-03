CREATE TABLE captcha_tasks (
    uuid VARCHAR NOT NULL,
    index BIGINT NOT NULL,
    difficulty  SMALLINT NOT NULL,
    math VARCHAR NOT NULL,
    answer INTEGER NOT NULL
);

ALTER TABLE captcha_tasks ADD CONSTRAINT captcha_tasks_pkey PRIMARY KEY (uuid);

CREATE INDEX index_idx ON captcha_tasks USING btree(index);