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