# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and I originally developed the initial version of this project as a part of my thesis, later decided to some extent rewrite the project, and included this in my very first article. This was originally my first project written in `Go`.

Feel free to copy parts of the code, fork the project, or use it as a reference. I am happy to receive feedback, and I am open to collaboration!

***

## About

This is `repository-metadata-aggregator`, which is an abstraction of parts of GitHub REST API, please refer to `docs/openapi.yaml` for the API docs. 

***

## How-To

### Run

- You require `go` and `make` to run this project.
- Setup environment variables (See: `.env.example`) and execute `make run`.

### cURL Commands

#### /api/v1/repositories/search/firstCreationDate

```bash
curl -X 'GET' \
  'http://localhost:8080/api/v1/repositories/search/firstCreationDate?query=language=Go&stars=>100' \
   --header 'Authorization: Bearer $GITHUB_TOKEN' 
```

#### /api/v1/repositories/search/lastCreationDate

```bash
curl -X 'GET' \
  'http://localhost:8080/api/v1/repositories/search/lastCreationDate?query=language=Go&stars=>100' \
   --header 'Authorization Bearer $GITHUB_TOKEN'
```

#### /api/v1/repositories/search

```bash
curl -X 'GET' \
  'http://localhost:8080/api/v1/repositories/search?firstCreationDate=2013-05-01&lastCreationDate=2013-05-01&language=Go&stars=>100' \
   --header 'Authorization Bearer $GITHUB_TOKEN'
```

***
