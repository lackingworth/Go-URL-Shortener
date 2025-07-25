# Go-URL-Shortener

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
