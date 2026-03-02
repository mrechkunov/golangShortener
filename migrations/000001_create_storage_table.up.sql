-- migrations/000001_create_storage_table.up.sql
-- Создание таблицы ссылок
CREATE TABLE storage (
    shortURL VARVHAR(255) PRIMARY KEY,
    originalURL VARCHAR(255) NOT NULL,
    uuid INTEGER NOT NULL
);

-- Базовый индекс для поиска по названию
CREATE INDEX idx_shortURL ON storage(shortURL);