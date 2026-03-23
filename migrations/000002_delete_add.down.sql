-- migrations/000002_delete_add.down.sql
-- Удаление столбца для удаленных ссылок
ALTER TABLE storage
DROP COLUMN IF EXIST isDeleted CASCADE;