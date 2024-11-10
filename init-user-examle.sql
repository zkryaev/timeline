-- Создаём пользователя с именем и паролем, которые берём из переменных окружения
CREATE USER king WITH PASSWORD 'kingpasswd';

-- Создаём базу данных, если она ещё не создана
CREATE DATABASE kingdom;

-- Назначаем владельцем базы данных нашего пользователя
ALTER DATABASE kingdom OWNER TO king;

-- Даем все права на базу данных нашему пользователю
GRANT ALL PRIVILEGES ON DATABASE kingdom TO king;