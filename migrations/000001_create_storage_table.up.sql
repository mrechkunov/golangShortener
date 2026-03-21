-- migrations/000001_create_storage_table.up.sql
-- Создание таблицы ссылок
CREATE TABLE IF NOT EXISTS storage (
    shortURL VARCHAR(255) PRIMARY KEY,
    originalURL VARCHAR(255) NOT NULL,
    count INTEGER NOT NULL,
    uuid INTEGER NOT NULL,
    cookie VARCHAR(255) NOT NULL,
);

-- Базовый индекс для поиска по названию
CREATE INDEX idx_shortURL ON storage(shortURL);