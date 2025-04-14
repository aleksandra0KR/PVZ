# Индексы:

- email, password таблицы users

Чтобы быстрее производить логирование пользователя в систему. Поскольку email стоит первым, то запрос из регистрации пользователя для проверки существования пользователя с таким email тоже будет выполняться быстрее


- reception_id, date_time таблица products

Для поиска последнего товара для операции удаления

- pvz_id, status таблица receptions

Для поиска последней активной приемки

## Схема базы данных
![](https://github.com/aleksandra0KR/PVZ/blob/master/scheme.png)


## Заметки про спецификацию OpenApi
В задании были указаны роли client, moderator, а в api спецификации employee, moderator. Были выбраны роли из спецификации.
В моей реализации у get /pvz есть response 401, поскольку по заданию только авторизованные пользователи могут иметь доступ к этой информации
Так же в спецификации в возвращающемся теле для /pvz/{pvzId}/close_last_reception
"status": "in_progress", но более логично выдавать "status": "close" чтобы отражать текущее состояние объекта и не вводить в заблуждение

# Запуск в Docker

Склонировать проект с гита

```
git clone https://github.com/aleksandra0KR/PVZ
```

Перейти в директорию проекта

```
cd PVZ
```

Забилдить

```
docker compose build
```

Запустить:

```
docker compose up
```

---

# Запустить без Docker

```
git clone https://github.com/aleksandra0KR/PVZ
```

Перейти в директорию проекта

```
cd PVZ
```

Запустить

```
go run cmd/main.go
```

### В файле .env можно поменять на нужные вам параметры


## Тесты
Покрытие без grpc файлов составляет 84.0%
``` 
go test -coverprofile=coverage.out -coverpkg=$(go list ./... | grep -v "grpcPVZ" | paste -sd, -) ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html 
```

Есть интеграционный тест:
- Получает нужные токены
- Создает новый ПВЗ
- Добавляет новую приёмку заказов
- Добавляет 50 товаров в рамках текущей приёмки заказов
- Закрывает приёмку заказов


### Сервисы:
- http по умолчанию на порту 8080
- grpc по умолчанию на порту 3000
- prometheus по умолчанию на порту 9000

