# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and I am building a research tool for my thesis. This is a complete open-source project, which started from my interest to learn `go`. Feel free to `fork`, leave comment, mail me, or copy parts of the code to your own usage. This `README.md` -file contains notes for the project, tasks, and descriptive information of the logic behind the functionality.

---

## About

This is `glsgen` (gls Generator) research tool which is a REST API, which collects data of open-source repositories and generates a `.csv` file from them. Tool uses `SourceGraph GraphQL API` and `GitHub GraphQL API`, and counts the sizes of repositories and their dependencies. Tools initial purpose is to create dataset for my thesis.

---

## Plugins

- `glsgen` is designed to be modular, `pkg/plugins` folder represents what kind of repositories can be analyzed. I am working (at the moment of writing, 27.10.2022), on `goplg`, which aims to offer functionality to analyze the quality of repositories, which primary language is `go`. Alpha Version of the Plugin is Completed in 15.1.2023.

- Go (Alpha Version Released)
- Planning: node.js

---

# Development

## How-To: Run

- See `Makefile`
- Requires: `go`, `postgresql`
- `.env` -file.

```
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_HOST=
POSTGRES_PORT=
GITHUB_USERNAME=
GITHUB_TOKEN=
GITHUB_GRAPHQL_API=https://api.github.com/graphql
SOURCEGRAPH_GRAPHQL_API=https://sourcegraph.com/.api/graphql
GOPATH=$HOME/go
GOPATH_PROD=/go
PROCESS_DIR=/path/to/workdir/tmp
PROCESS_DIR_PROD=/path/to/workdir/tmp
MAX_GOROUTINES=64
```

### Generate Metadata

- Endpoint `/api/gls/v1/repository/fetch/metadata` - requires following query parameters `type` and `count`.
    - Example: `/api/gls/v1/repository/fetch/metadata?type=go&count=1`

---

## How-To: Contribute

- I don't have a structured way to accept contributions to the project, but feel free to leave a `pull request`, if you feel like it. :)

---

# Notes

- Implementation might have issues, since I am implementing concurrency for the first time.

## TODO

- TBD

## TODO (Out of the Thesis Scope)

- Quality Measure -endpoint.
- .env -> config.yaml
- Commit -analysis.
- GitLab Runner.

---

## Docker Notes

### Run gls a a Docker Container

- Network:`docker network create --subnet=172.19.0.0/16 <network_name>`
- Build: `docker build --tag gls:latest .`
- Run: `docker run -idt -p 8080:8080 --cpus=<cpu_count> -m <memory_amount> --name gls --net <network_name> gls:latest`

You might need to configure `.env` - file values to fit your environment.

### Run PostgreSQL as a Docker Container

- Create Docker Network: `docker network create --subnet=172.19.0.0/16 <network_name>`
- Run the PostgreSQL container in a certain network, with couple environment variables, and static IP -address.

```docker run -d --name postgres --net <network_name> --ip 172.19.0.2 -e POSTGRES_USER=<username> -e POSTGRES_PASSWORD=<password> -e POSTGRES_DB=<database_name> -p 5432:5432 postgres```

- Verify, that the PostgreSQL is running, with: `psql -h localhost -U postgres`

### General Docker Notes

- Check the IP -address of the container: `docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' container_name`

---
