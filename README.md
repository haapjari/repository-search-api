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

### Build

- Build the Image: `docker build -t repository-search-api:latest .`
- Run the Image (On the Host, for Simplicity): `docker run -idt -e PORT=8080 --network=host repository-search-api:latest`

### Invoke

#### /api/v1/repos/search

```bash
curl "localhost:8080/api/v1/repos/search?firstCreationDate=2008-01-01&lastCreationDate=2009-01-01&language=Go&minStars=100&maxStars=1000&order=desc" --header "Authorization: Bearer $GITHUB_TOKEN"
```

---