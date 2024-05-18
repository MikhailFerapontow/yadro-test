# yadro-test
Тестовое задание для вакансии Инженер по разработке ПО для базовых станций (Go) в рамках Стажировка "Импульс" YADRO 2024.

Участник: Ферапонтов Михаил Владимирович

# Требования

- Makefile
- Docker

# Инструкция по запуску

### Создание образа
```sh
docker build -t "image_name" .
```
или
```sh
make build #создаст образ с именем ferapontov/yadro-test
```

### Запуск программы в контейнере
```sh
docker run --rm \
	--mount type=bind,source="/path/to/input/file",target=/app/input.txt \
	-w /app \
	"image_name" ./app input.txt
```
или
```sh
make run INPUT="/path/to/input/file"
```

### Запуск тестов
```
docker run --rm "image_name" go test ./...
```
или
```sh
make test
```

# Примечание
При событии 3 непонятно что должно происходить, если клиент не находится в клубе.

Было принято решение, в такой ситуации генерировать ошибку "ClientUnknown"