### Архитектура
 <div class="white-background">
      <img src="https://github.com/user-attachments/assets/e447d923-744a-4da1-860a-c2d422631bfc" alt="Описание изображения">
 </div>

### Конвертация часовых поясов
> Предварительно при запуске объект `loader` проходит по указанному URL в .env и получает данные, затем грузит их в память.</br>

![image](https://github.com/user-attachments/assets/0321248a-de21-453f-bbf1-e749ead620f3)

### Авторизация
> На каждый хендлер задаются права доступа как в UNIX системах. Нельзя обратиться к чужим данным поскольку в запросах к БД исплользуется id сущности, который извлекается из полученного токена, кроме тех мест где это разрешено через query-параметры для получения открытой информации</br>

**Шаблон**</br>

![image](https://github.com/user-attachments/assets/5732ab3c-fe62-4ff3-bbb4-c049ef6adedc)

**Пример**</br>
![image](https://github.com/user-attachments/assets/43112147-0b3b-4838-8fad-e3ce87bc5af8)

### Схема мастер БД
![image](https://github.com/user-attachments/assets/c19e8e91-a1a1-44c3-ba7a-85b148f3b034)

### Мониторинг запроса
> \*Метрики обновляются в мидлваре, но черпаются с отдельной ручки, отдельного роутера.</br>

![image](https://github.com/user-attachments/assets/99c455a5-10b2-4824-be98-8963ecae9ccd)

### Работа с почтой
> \*Письма при падении внешнего сервиса теряются, решается добавлением брокера сообщений с соответствующими настройками</br>

![image](https://github.com/user-attachments/assets/ff490154-e64f-437a-80c0-ed650ec5be33)

### Как храним изображения в S3
**Для главных сущностей**</br>
![image](https://github.com/user-attachments/assets/7d2a3ef3-902f-412d-8e9b-84c050ba589b)</br>
**Для суб-сущностей (работники, изображения галереи)**</br>
![image](https://github.com/user-attachments/assets/5bf79707-7c77-4f1a-86cb-66a3c4e9abbb)
