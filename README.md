# Тестовое задание

Необходимо спроектировать и разработать сервис, 
который будет предоставлять API для работы с информацией о погоде.

## Часть 1. Получение данных.

Используя [geocoding-api](https://openweathermap.org/api/geocoding-api) получить информацию о любых 20-ти городах и сохранить ее в локальную базу. Необходимо сохранить название, страну, и координаты (понадобятся для получения погоды). Сохранять данные о городах можно любым удобном для вас способом.

Используя открытый API [open weather map](https://openweathermap.org/forecast5) и координаты городов получить предсказание погоды на 5 дней и сохранить результат в БД.

При запуске сервис должен запросить предсказания погоды из внешнего источника и сохранить их в локальную базу данных. Если предсказания для конкретного города и времени в БД уже есть необходимо обновить существующую запись.

Нет необходимости хранить все поля ответа в отдельных колонках БД. В отдельных колонках должны лежать:

- `temp` — температура днем
- `date` — дата предсказания
- Всю остальную информацию нужно хранить в формате json

```
💪 Если задание кажется вам простым.
Дополнительно. Сделать фоновый процесс который будет раз в минут асинхронно обновлять данные о погоде.
```

```
💪 Если все еще слишком просто.
Дополнительно. Распараллелить процесс получения информации о погоде из внешнего API.
```

В проект необходимо добавить файлик migrate.sql который должен содержать в себе sql-запросы для создания схемы вашей БД со всем таблицами, индексами и т.д.

```
💪 Для тех кто знает докер.
Оберните ваш проект в docker. Сервис должен запускаться через doсker-compose и после поднятия контейнеров полностью сконфигурирован и готово к работе.
```

## Часть 2.  Разработка API.

Какой функционал нужен пользователям API:

1. Список городов, для которых есть предсказания о погоде (отсортированный по названию). Из данного списка пользователь сможет выбрать город для просмотра полной информации по нему.
2. Список с кратким предсказанием для выбранного города. Ответ должен содержать: страну, название города, среднюю температуру на весь доступный **будущий** период, список дат для которых доступно предсказание. Дата должна быть отсортирована в хронологическом порядке. Пользователь сможет выбрать дату для получения полной информации на этот день.
3. Детальная информация о погоде для конкретного города и конкретного времени. Ответ должен быть максимально полным и содержать в себе всю информацию, которую получилось достать из внешнего API на выбранное время.

```
💪 Для самых смелых.
Дополнительно. Добавить в сервис функционал пользователей и спроектировать API для добавление города в список избранных пользователя. Жестких требований к этой части задания нет, оставляем пространство для творчества 🙂
```

## Важные комментарии

- Сервис должен быть написан на Golang, запрещено использовать ORM для работы с БД. Все запросы должны быть сделаны на чистом SQL. В остальном ограничений по библиотекам и фреймворкам нет.
- База данных — PostgreSQL.
- В первую очередь необходимо реализовать основной функционал, затем из секции «дополнительно»
- Качество важнее скорости. На выполнение задачи дается неделя. Если остается время рекомендую потратить его на улучшение вашего решения. Напишите документацию, unit-тесты, добавьте логгирование и т.д.
- Не стесняйтесь задавать вопросы, если что-то не понятно в условии задания. Лучше лишний раз спросить чем сделать задачу не верно 😉.

## Installation

Запустить сервис командой `make build`

API:
1. Список городов, для которых есть предсказания о погоде `Get /cities`
2. Список с кратким предсказанием для выбранного города `Get /cities/{city}`
3. Детальная информация о погоде для конкретного города и конкретного времени `Get /cities/{city}?date={yyyy-MM-dd HH:mm:ss}`
