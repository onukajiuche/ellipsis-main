# brief
Brief makes your url's as `brief` as the word: `brief` (with a few added chars 🤪)

Brief uses MD5 hashing algorithm, with base64 URL encoding and a unique counter to
achieve unique short urls for each long url. Brief also allows you to specify a specific 
url of your choice. 

## Features
- [x] User Authentication
- [x] Hashing Algorithm (MD5 + base64 URL Encoding + Random counter)
- [x] API Testing
- [ ] Caching with Redis
- [ ] Analytics
- [x] Swagger Documentation
- [x] Postman Documentation
- [ ] Continuos Integration/Delivery (CI/CD)

## Requires
`go 1.17+` `postgresql` `redis` `docker (optional)`

## Setup

- Clone project using:
```bash
    git clone https://github.com/emmrys-jay/brief.git
```

- Create your env file named `mine.env`, and write your env variables according to
`sample.env` in the projects root directory

- Note: If you have docker installed, you can use postgres from docker by running the 
following in a separate bash shell.
```bash
    docker-compose -f postgres.docker-compose.yml up
```

- Run unit tests using:
```bash
    go test -tags=unit ./...
```

- Ensure your database is up and configured, then run unit tests using:
```bash
    go test -tags=integration ./...
```

- Start server using:
```bash
    go run main.go
```

## Documentation

- After starting the server, you can view swagger documentation at
[http://localhost:8080/swagger/index.html#](http://localhost:8080/swagger/index.html#)

- View postman documentation for the API at 
[Postman link](https://documenter.getpostman.com/view/20046026/2s946feYbR)
