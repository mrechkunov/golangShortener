
-- migrations/000002_delete_add.up.sql
-- Добавление столбца для удаленных ссылок
ALTER TABLE storage ADD COLUMN isDeleted BOOLEAN DEFAULT false;