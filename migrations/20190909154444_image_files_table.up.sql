CREATE TABLE images
(
  id           BIGSERIAL PRIMARY KEY,
  source_name   TEXT NOT NULL,
  stored_name   TEXT NOT NULL
);

COMMENT ON TABLE images IS 'Загруженные изображения';
COMMENT ON COLUMN images.source_name IS 'Исходное имя файла';
COMMENT ON COLUMN images.stored_name IS 'Сохранненное имя файла';