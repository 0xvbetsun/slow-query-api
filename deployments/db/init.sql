CREATE EXTENSION IF NOT EXISTS pg_stat_statements SCHEMA PUBLIC CASCADE;

CREATE TABLE DATA (
    id serial NOT NULL PRIMARY KEY,
    name varchar(20) NOT NULL,
    duration int NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

INSERT INTO
    data(name, duration)
VALUES
    ('first', 20),
    ('second', 526),
    ('first_slow', 4214);

UPDATE
    DATA
SET
    duration = (
        SELECT
            duration + 30
        FROM
            pg_sleep(1)
    )
WHERE
    value < 200;

DELETE FROM
    DATA
WHERE
    duration > (
        SELECT
            100
        FROM
            pg_sleep(1)
    );

SELECT
    sum(duration),
    pg_sleep(3)
FROM
    DATA;

SELECT
    *
FROM
    pg_stat_statements;