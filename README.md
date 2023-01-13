# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and I am building a research tool for my thesis. This is a complete open-source project, which started from my interest to learn `go`. Feel free to `fork`, leave comment, mail me, or copy parts of the code to your own usage. This `README.md` -file contains notes for the project, tasks, and descriptive information of the logic behind the functionality.

---

## About

This is **Glass**, a research tool which aims to offer data collection capabilities to measure quality of open-source repositories and return a single value called "Quality Measure" to represent state of repositories. **Glass** is essentially an API, that collects data from multiple datasources, such as SourceGraph GraphQL API, GitHub GraphQL API and GitHub REST API, and combines that to meaningful form, that can be analyzed. 

**Glass** is going to be used to create a dataset for my thesis, and could be further developed to a tool, that can run inside GitHub Actions or GitLab CI/CD as a separate tool.

---

## Plugins

- **Glass** is designed to be modular, `pkg/plugins` folder represents what kind of repositories can be analyzed. I am working (at the moment of writing, 27.10.2022), on `goplg`, which aims to offer functionality to analyze the quality of repositories, which primary language is `go`.

- *WIP*: Go
- Proposed: node

---

# Development

## How-To: Run

- See `Makefile`
- Requires: `go`, `postgresql`
- `.env` -file, you need to fill up these values: <!-- TODO: Theres multiple hardcoded values, give these examples to here.>

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
BASEURL=localhost:8080/api/glass/v1
GOPATH=
PROCESS_DIR=
MAX_GOROUTINES=
```

---

## How-To: Contribute

- I don't have a structured way to accept contributions to the project, but feel free to leave a `pull request`, if you feel like it. :)

---

# Notes

## TODO

- Optimizations
- Performance Testing
- Add the skipped libraries on the data structures, and try to process them again.
- Prune other than "Go" Primary Language repositories, before they are processed.

---

## Performance Testing

**13.1.2023**

- Benchmark, 1 repository: 165 sec
- Benchmark, 10 repositories: Crash

**14.1.2023**

- Benchmark, 1 repository: 175 sec, 118 sec, 90 sec (64 routines), 85 sec (128 Routines)
- Benchmark

---

## Data Collection and Research

### Primary Repositories

- SourceGraph API, Query: `lang:go select:repo repohasfile:go.mod count:100000`
    - Repository Name
    - Repository URL

```
query {
  search(query: "lang:go AND select:repo AND repohasfile:go.mod AND count:100000", version: V2) {
    results {
        repositories {
            name
        }
    }
  }
}
```

- GitHub GraphQL API: Following Query Returns
    - Total Count of Commits in Repository
    - Total Count of Open Issues in Repository
    - Total Count of Closed Issues in Repository
    - Total Size of the Repository
    - Creation Date
    - Primary Language
    - Stars Count
    - License Info
    - Latest Release Date

```
query {
        repository(owner: "TBD", name: "TBD") {
            createdAt
            primaryLanguage {
                name
            }
            defaultBranchRef {
                target {
                    ... on Commit {
                        history {
                            totalCount
                        }
                    }
                }
            }
            openIssues: issues(states:OPEN) {
                totalCount
            }
            closedIssues: issues(states:CLOSED) {
                totalCount
            }
            languages {
                totalSize
            }
            stargazerCount
            licenseInfo {
                key            
            }
            latestRelease {
                name
                publishedAt
            }         
        }
}
```

---

### Quality Measure

- Repository Activity: Higher -> Better
    - Amount of commits.
- Maintainers: Higher -> Better
    - Amount of maintainers. 
- Ratio of Open Issues to Closed Issues: Less -> Better
    - Amount of Open Issues
    - Amount of Closed Issues
- Creation Date: Older -> Better
    - Might be an indicator of maturity of the repository.
- Stars: Higher -> Better
    - Determines the popularity of the repository.
- Releases: Higher -> Better
    - Determines the maturity of the repository, more releases might indicate more mature project.
- Latest Release Date: More Recent -> Better
    - If there are more than certain threshold amount of time from last release, might be worse.

Thresholds of these amounts will be calculated, thresholds will be inbeween 0-5, where 2.5 is at middle of the amounts.

These values will be averaged in a single `Quality Measure`. Correlation will be calculated ratio of library to original code lines, or ratio of sizes. Is there a correlation between bigger ratio and quality measure.

#### Derivative Information

- Correlation:
    - QM / Original Codebase Size
    - QM / Ratio of (Open Issues / Closed Issues)
    - QM / Maintainers
    - QM / Creation Date
    - QM / Stars

---

### Database Tables

- Table: "Repository"
    - Primary Key: RepositoryId
    - Repository Struct: Repository Name, Url, CommitCount, Collaborators, Open Issues, Closed Issues, Original Codebase Size, Total Library Codebase Size, ProjectType, PrimaryLanguage
- Table: "Commits"
    - Primary Key: CommitId
    - Columns: RepositoryId, Commit Date, Commit User, Repository Name

---

### Run PostgreSQL as a Docker Container

- Create Docker Network: `docker network create --subnet=172.19.0.0/16 glass_network`
- Run the PostgreSQL container in a certain network, with couple environment variables, and static IP -address.

```docker run -d --name postgres --net glass_network --ip 172.19.0.2 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -p 5432:5432 postgres```

- Verify, that the PostgreSQL is running, with: `psql -h localhost -U postgres`

---
