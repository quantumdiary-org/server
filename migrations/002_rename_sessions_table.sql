-- +goose Up
-- +goose StatementBegin
-- Переименовываем таблицу, если она существует
ALTER TABLE net_school_sessions RENAME TO sessions_backup;
-- Создаем новую таблицу с правильным именем
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL UNIQUE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    netschool_url TEXT NOT NULL DEFAULT 'https://sgo.rso23.ru',
    school_id INTEGER NOT NULL,
    student_id TEXT NOT NULL,
    year_id TEXT NOT NULL,
    api_type TEXT NOT NULL DEFAULT 'ns-webapi',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);
-- Копируем данные из резервной таблицы
INSERT INTO sessions SELECT * FROM sessions_backup;
-- Удаляем резервную таблицу
DROP TABLE sessions_backup;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Переименовываем таблицу обратно
ALTER TABLE sessions RENAME TO net_school_sessions_backup;
-- Создаем старую таблицу
CREATE TABLE net_school_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL UNIQUE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    netschool_url TEXT NOT NULL DEFAULT 'https://sgo.rso23.ru',
    school_id INTEGER NOT NULL,
    student_id TEXT NOT NULL,
    year_id TEXT NOT NULL,
    api_type TEXT NOT NULL DEFAULT 'ns-webapi',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);
-- Копируем данные обратно
INSERT INTO net_school_sessions SELECT * FROM net_school_sessions_backup;
-- Удаляем временную таблицу
DROP TABLE net_school_sessions_backup;
-- +goose StatementEnd