# Go-URL-Shortener-For-Ozon-Test
___
#### Русский:
___
## Установка
* У вас должен быть установлен [Go v1.21.5](https://go.dev/doc/install) (или выше)
* У вас должна быть установлена актуальная версия [Docker](https://www.docker.com/) 

## Запуск контейнеров / программы

* Клонируйте данный репозиторий в выбранную вами локацию
* *(Необязательно)* Измените предоставленные ```.env``` и ```Dockerfile``` файлы для кастомизации сетевых подключений
* Откройте консоль разработчика / Bash / PowerShell
* Зайдите в корневую директорию склонированного репозитория (там, где находится файл ```docker-compose.yml```)
  с помощью ```cd "путь/к/папке"```
> [!IMPORTANT]    
> Для запуска api с использованием in-memory хранилиша (Redis) используйте:
> ```
> docker-compose run -d --service-ports --name "api" api -memory
> ```
>
> Для запуска api с использованием Postgres SQL используйте:
> ```
> docker-compose run -d --service-ports --name "api" api -db
> ```

> [!TIP]
> Флаги и пояснения
>
> ```docker compose run``` - запускает предоставленный в команде и ```docker-compose.yml``` сервис
> 
> ```-d``` - запускает контейнер на заднем плане, позволяя продолжать использование терминала
> 
> ```--service-ports``` - создает и отлаживает порты предоставленные в ```Dockerfile``` файлах и ```docker-compose.yml```
> 
> ```--name "имя для контейнера"``` - предоставляет контейнеру имя
> 
> ```api``` - имя сервиса, указанное в ```docker-compose.yml```
> 
> ```-arg``` - кастомный аргумент, который определяет какое хранилище следует использовать в программе. Возможные варианты - **memory** и **db** для использования **Redis** и **Postgres** соответственно.
>
>  Отсутствие этого флага запустит сервис по стандартной конфигурации - с использованием Redis

## Использование (Запросы)

Имеются 2 API эндпоинта:
> [!NOTE]  
> 
> ```*DOMAIN*/api``` - принимает ссылку и возвращает ее сокращенный вариант - стандартный эндпоинт - ```localhost:3000/api``` - POST запрос
> 
> ```*DOMAIN*/:url``` - возвращает JSON-объект с оригинальной ссылкой - стандартный эндпоинт - ```localhost:3000/:url``` - GET запрос

> [!IMPORTANT] 
> Чтобы поменять функционал JSON ответа на автоматическое перенаправление необходимо заменить последнюю строку в файлах ```routesPostgres.go``` и ```routesRedis.go``` на комментарий над ней

> [!TIP]
> * Чтобы получить сокращенную ссылку, отправьте запрос POST на адрес ```*DOMAIN*/api```, в теле запроса укажите ссылку, которую нужно сократить в формате JSON:
>
> ```{"url":"ваша ссылка"}```
>
> * Чтобы получить кастомную "красивую" сокращенную ссылку, предоставьте ее вместе с обычным url в теле POST запроса:
>
> ```{"url":"ваша ссылка", "short":"ваша красивая ссылка"}```

## Особенности

* Использование двух разных хранилищ (Redis и Postgres SQL), с возможностью их поменять и кастомизировать
* Настраиваемый Rate limiter для частых запросов api (<20 запросов за 10 секунд)
* Настраиваемый Rate limiter для in-memory хранилища
* Возможность получить "красивую" ссылку
* Стандартные сокращенные ссылки состоят из 10 символов (настраивается) - символов латинского алфавита в нижнем и верхнем регистрах, цифрах и граунда ("_")
* Автоматический рестарт хранилищ и api при неполадках

>[!NOTE]
> Протестировано с помощью Unit-тестов и Postman

___
#### English:
___

## Installing
* You must have [Go v1.21.5](https://go.dev/doc/install) (or higher) installed on your system
* You must have up-to-date version of [Docker](https://www.docker.com/) installed on your system 

## Starting containers / executing program

* Clone this repository to the location of your choosing
* *(Optional)* Change provided ```.env``` and ```Dockerfile``` files for networking customization
* Open Bash / PowerShell terminal
* Navigate to the root directory of the cloned repository (where ```docker-compose.yml``` is located)
  using ```cd "path/to/directory"```
> [!IMPORTANT]    
> To start api with in-memory storage (Redis) run the following command:
> ```
> docker-compose run -d --service-ports --name "api" api -memory
> ```
>
> To start api with Postgres SQL run the following command::
> ```
> docker-compose run -d --service-ports --name "api" api -db
> ```

> [!TIP]
> Flags and explanation
>
> ```docker compose run``` - starts up the service provided in the command line and in ```docker-compose.yml```
> 
> ```-d``` - runs container in background, allowing to continue using your terminal
> 
> ```--service-ports``` - enables and maps to host ports provided in ```Dockerfile``` files and ```docker-compose.yml```
> 
> ```--name "custom name"``` - naming the container
> 
> ```api``` - service name, provided in ```docker-compose.yml```
> 
> ```-arg``` - custom argument which decides what type of storage to implement. Possible variants are - **memory** and **db** for using **Redis** and **Postgres** respectively.
>
> If this flag is not provided - service runs with default configuration - using Redis

## Usage (Requests)

There are 2 API endpoints:
> [!NOTE]  
> 
> ```*DOMAIN*/api``` - accepts url and returns its shortened variant - default endpoint - ```localhost:3000/api```
> 
> ```*DOMAIN*/:url``` - returns JSON object with the original url - default endpoint - ```localhost:3000/:url```

>[!IMPORTANT]  
> To change JSON response to automatic redirect functionality, you need to swap the last line of code in ```routesPostgres.go``` and ```routesRedis.go``` files with the comment above it

> [!TIP]
> * To get short url, send POST request to ```*DOMAIN*/api``` and provide in the body JSON-formatted original url you wish to shorten: 
>
> ```{"url":"your url"}```
>
> * To get custom "beautiful" short url, provide it with original full url in the body of the POST request:
>
> ```{"url":"your_url", "short":"your_custom_short"}```

## Features

* Customizable, changable dual storage capability (Redis and Postgres SQL)
* Customizable Rate limiter for rapid api requests (<20 requests in 10 seconds)
* Customizable Rate limiter for in-memory storage
* Ability yo provide custom short url
* Default short urls are 10 characters long (customizable) - [a-z] and [A-Z] characters, numbers and low dash ("_")
* Automatic storage and api restart if malfunction occurred

## Version History

* v.0.0.1:

    * Initial Release
