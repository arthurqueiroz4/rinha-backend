CREATE TABLE consumers (
        id SERIAL PRIMARY KEY,
        bound INTEGER NOT NULL,
        balance INTEGER NOT NULL
);

CREATE TABLE transactions (
        id SERIAL PRIMARY KEY,
        consumer_id INTEGER NOT NULL,
        type VARCHAR(1) NOT NULL,
        value INTEGER NOT NULL,
        description VARCHAR(10) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
INSERT INTO consumers (bound, balance)
VALUES
    (1000 * 100, 0),
    (800 * 100, 0),
    (10000 * 100, 0),
    (100000 * 100, 0),
    (5000 * 100, 0);
END; $$