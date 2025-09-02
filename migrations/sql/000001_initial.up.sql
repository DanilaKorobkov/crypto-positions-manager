BEGIN;

CREATE TABLE subjects (
    telegram_user_id INTEGER PRIMARY KEY,
    payload JSONB
);

COMMIT;