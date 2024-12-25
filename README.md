# Описание

Документация к проекту появится до конца декабря...

# Запуск и прочее
- `task deploy` - проверит ENVs, развернет докер, произведет миграции, запустит приложение. После завершения работы приложения, завершит работу докера, но контейнеры останутся!. <br>
- `task docker-rm` - удалить существующие контейнеры <br>
- `task build, task run` - билд и запуск конкретно приложения <br>
- `task swagger` - сгенерирует Swagger-документацию <br>

### Используемые библиотеки
[Роутер: gorilla/mux](https://github.com/gorilla/mux) <br>
[Миграции: golang-migrate](https://github.com/golang-migrate/migrate?tab=readme-ov-file) <br>
[Работа с JWT: golang-jwt](https://github.com/golang-jwt/jwt) <br>
[Работа с JSON: json-iterator](https://github.com/json-iterator/go) <br>
[Драйвер для БД: pgx](https://github.com/jackc/pgx) <br>
[Интерфейс для БД: sqlx](https://github.com/jmoiron/sqlx) <br>
[Работа с почтой по SMTP: gomail](https://github.com/go-gomail/gomail) <br>
[Валидация полей структур: validator](https://github.com/go-playground/validator) <br>
[Работа с переменными окружения: godotenv](https://github.com/joho/godotenv) <br>
[Генерация Swagger документаци: swaggo](https://github.com/swaggo/swag) <br>
[Фоновые задачи через CRON: gocron](https://pkg.go.dev/github.com/go-co-op/gocron/v2) <br>
[S3 MinIO SDK: go-minio](https://github.com/minio/minio-go) <br>





