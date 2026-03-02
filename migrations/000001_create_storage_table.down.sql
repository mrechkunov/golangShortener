-- migrations/000001_create_movies_table.down.sql
-- Откат создания таблицы фильмов
DROP INDEX IF EXISTS idx_shortURL;
DROP TABLE IF EXISTS storage;