# yadro-test
Тестовое задание для вакансии Инженер по разработке ПО для базовых станций (Go) в рамках стажировки "Импульс" YADRO 2024.

Участник: Ферапонтов Михаил Владимирович

## Требования

- Makefile
- Docker

## Инструкция по запуску

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

### Удаление образа
```sh
make clean # удалит образ ferapontov/yadro-test
```

## Тестовые примеры

Все примеры находятся в папке /examples.

[000_valid.txt](examples/000_valid.txt): Пример из задания.

[001_empty.txt](examples/001_empty.txt): Пустой файл

Ожидаемый результат:
```sh
line 1 invalid number of tables:
```

[002_invalid_hours.txt](examples/002_invalid_hours.txt): Неправильно указаны часы работы клуба.

Ожидаемый результат:
```sh
line 2 invalid work hours:
asdfasdfadsf
```

[003_invalid_tariff.txt](examples/003_invalid_tariff.txt): Неправильно указано значение почасовой оплаты.

Ожидаемый результат:
```sh
line 3 invalid hour payment:
sdfgsdfg
```

[004_invalid_command.txt](examples/004_invalid_command.txt): Неправильно указана команда

Ожидаемый результат:
```sh
line 5 invalid command:
as;dfja;sdf;adsf;
```

[005_invalid_time.txt](examples/005_invalid_time.txt): Неправильно указано время команды. Время предыдущей команды больше времени текущей.

Ожидаемый результат:
```sh
line 7 Time can only flow forward...
09:10 3 client1
```

[006_invalid_table.txt](examples/006_invalid_table.txt): Номер стола, за который хочет сесть клиент, указан неправильно.

Ожидаемый результат:
```sh
line 11 Table number out of range:
10:59 2 client3 123
```

[007_valid2.txt](examples/007_valid2.txt): Дополнительный правильный пример.

## Примечание
При событии 3 непонятно, что должно происходить, если клиент не находится в клубе.

Было принято решение, в такой ситуации генерировать ошибку "ClientUnknown"
