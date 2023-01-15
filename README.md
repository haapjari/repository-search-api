# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and I am building a research tool for my thesis. This is a complete open-source project, which started from my interest to learn `go`. Feel free to `fork`, leave comment, mail me, or copy parts of the code to your own usage. This `README.md` -file contains notes for the project, tasks, and descriptive information of the logic behind the functionality.

---

## About

This is **Glass**, a research tool which aims to offer data collection capabilities to measure quality of open-source repositories and return a single value called "Quality Measure" to represent state of repositories. **Glass** is essentially an API, that collects data from multiple datasources, such as SourceGraph GraphQL API, GitHub GraphQL API and GitHub REST API, and combines that to meaningful form, that can be analyzed. 

**Glass** is going to be used to create a dataset for my thesis, and could be further developed to a tool, that can run inside GitHub Actions or GitLab CI/CD as a separate tool.

---

## Plugins

- **Glass** is designed to be modular, `pkg/plugins` folder represents what kind of repositories can be analyzed. I am working (at the moment of writing, 27.10.2022), on `goplg`, which aims to offer functionality to analyze the quality of repositories, which primary language is `go`.

- Go
- Proposed: node

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
GOPATH=
GOPATH_PROD=
PROCESS_DIR=
PROCESS_DIR_PROD=
MAX_GOROUTINES=64
```

### Generate Metadata

- Endpoint `/api/glass/v1/repository/fetch/metadata` - requires following query parameters `type` and `count`.
    - Example: `/api/glass/v1/repository/fetch/metadata?type=go&count=1`

---

## How-To: Contribute

- I don't have a structured way to accept contributions to the project, but feel free to leave a `pull request`, if you feel like it. :)

---

# Notes

- Implementation might have issues, since I am implementing concurrency for the first time.

## TODO

- 2023/01/15 04:48:18 parse "https://raw.githubusercontent.com/rudderlabs/rudder-server/master//\tWeshouldfrequentlyreviewthissectiontoremoveorupdatethe/go.mod": net/url: invalid control character in URL
- Error

## TODO (Out of the Thesis Scope)

- Quality Measure -endpoint.
- GitLab Runner.
- More precise QM -endpoint.

---

## Run PostgreSQL as a Docker Container

- Create Docker Network: `docker network create --subnet=172.19.0.0/16 glass_network`
- Run the PostgreSQL container in a certain network, with couple environment variables, and static IP -address.

```docker run -d --name postgres --net glass_network --ip 172.19.0.2 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -p 5432:5432 postgres```

- Verify, that the PostgreSQL is running, with: `psql -h localhost -U postgres`

- Check address of the container: 

```docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' glass_postgres_1```

---
