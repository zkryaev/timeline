# Описание

# Запуск и прочее
- `task deploy` - проверит ENVs развернет докер, произведет миграции, запустит приложение. После завершения работы приложения завершит работу докера. <br>
- `task run` - запуск конкретно приложения <br>
- `task gen-api` - сгенерирует Swagger-документацию <br>

### Схема БД
![image](https://github.com/user-attachments/assets/a9277e30-09dd-483d-b710-bc61f7ca2c70)

### Зависимости
[Роутер](https://github.com/gorilla/mux) <br>
[Миграции](https://github.com/golang-migrate/migrate?tab=readme-ov-file) <br>
[Работа с JWT](https://github.com/golang-jwt/jwt) <br>
[Работа с JSON](https://github.com/json-iterator/go) <br>
[Драйвер для БД](https://github.com/jackc/pgx?ysclid=m3kfe9usdw209993533) <br>
[Интерфейс для БД](https://github.com/jmoiron/sqlx) <br>
[Работа с почтой по SMTP](https://github.com/go-gomail/gomail) <br>
[Валидация полей структур](https://github.com/go-playground/validator) <br>
[Работа с переменными окружения](https://github.com/joho/godotenv) <br>
[Генерация Swagger документации](https://github.com/swaggo/swag) <br>






