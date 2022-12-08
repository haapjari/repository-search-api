# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and I am building a research tool for my thesis. This is a complete open-source project, which started from my interest to learn `go`. Feel free to `fork`, leave comment, mail me, or copy parts of the code to your own usage. This `README.md` -file contains notes for the project, tasks, and descriptive information of the logic behind the functionality.

---

## About

This is **Glass**, a research tool which aims to offer functionality to measure quality of open-source repositories and libraries and return a single value called "Quality Measure" to represent state of repositories. **Glass** is essentially an API, that collects data from multiple datasources, such as SourceGraph GraphQL API, GitHub GraphQL API and GitHub REST API, and combines that to meaningful form, that can be analyzed. 

**Glass** is going to be used to create a dataset for my thesis, where I am researching topic `Effects Software Reuse to Quality of Codebase`. 

---

## Plugins

- **Glass** is designed to be modular, `pkg/plugins` folder represents what kind of repositories can be analyzed. I am working (at the moment of writing, 27.10.2022), on `goplg`, which aims to offer functionality to analyze the quality of repositories, which primary language is `go`.

- *WIP*: Go

---

# Development

## How-To: Run

- See `Makefile`
- Requires: `go`, `postgresql`

---

## How-To: Contribute

- I don't have a structured way to accept contributions to the project, but feel free to leave a `pull request`, if you feel like it. :)

---

# Notes

## TODO

- Support for monorepo, this only supports singlerepos.
- docker-compose.yml
    - cAdvisor
    - Unable to mount environment variables to the container.
- Dockerfile 
    - Unable to Create Post Request from inside the container, getting `x509: certificate signed by unknown authority` error.
- Authorization
- Swagger API Documentation
- Logic
    - ~~Capability to fetch Repository Metadata~~
    - ~~Capability to fetch Primary Type Data~~
    - WIP: **Capability to fetch Library Type Data** -> Scoped out, there were so many corner cases, scoped this out.
    - Capability to fetch Commit Data on Repository Basis -> Scoped out, for time limitations.
    - Capability to fetch Code Smells (?) -> Scoped out, for time limitations.
    - GoReportCard Integration (?) -> Scoped out, for time limitations.

(?) - Will these be in the scope of thesis.

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
        }
}
```

### Libraries

- SourceGraph GraphQL API

```
query {
 repository(name: "github.com/TBD/TBD") {
    defaultBranch {
      target {
        commit {
          blob(path: "go.mod") {
            content
          }
        }
      }
    }
  }
}

```

- Following query returns the content of the "go.mod" - which can be used to parse the dependencies.

---

### Commits

<!-- TODO: Think, will this be implemented in the scope of the thesis. -->
- `script/list_commits.sh` 
    - Total Count of Commits in Repository
    - Every Commit in the Repository
    - Date of Commit in the Repository
    - Author of Commit in the Repository
- See: https://docs.github.com/en/rest/commits/commits

---

<!-- TODO: Think, will this be implemented in the scope of the thesis. -->
### Code Smells
- `gosec -fmt=json -out=results.json -stdout ./...`
    - Total Count of Code Smells in Repository

---

### Quality Measure

- Repository Activity: Higher -> Better
    - Amount of commits, with dates. 
- Maintenance: Higher -> Better
    - Collaborators  
- Code Smells (https://github.com/securego/gosec): Less -> Better
    - Total Count of code smells in the repository.
        - Code smells severity average.
- Ratio of Open Issues to Closed Issues: Less -> Better
    - Amount of Open Issues
    - Amount of Closed Issues
- Thresholds of these amounts will be calculated, threshholds will be inbeween 0-5, where 2.5 is at middle of the amounts.
- These values will be averaged in a single QM value. Correlation will be calculated ratio of library to original code lines, or ratio of sizes. Is there a correlation between bigger ratio and quality measure.
- Stars: Higher -> Better
    - Determines the popularity of the repository.
    - TODO: Refactor to use "Imports", when `https://pkg.go.dev/` has an API. (Follow: https://github.com/golang/go/issues/36785)
- Creation Date: Older -> Better
    - Might be an indicator of maturity of the repository.

#### Derivative Information

- Ratio: Library Codebase Size / Original Codebase Size
- Ratio: Open Issues / Closed Issues
- Commit Activity
- Maintentance: Collaborator Activity
- Which size packages are used in ratio (?)
- What license does repository use.

---

### Database Tables

- Table: "Repository"
    - Primary Key: RepositoryId
    - Repository Struct: Repository Name, Url, CommitCount, Collaborators, Open Issues, Closed Issues, Original Codebase Size, Total Library Codebase Size, ProjectType, PrimaryLanguage
- Table: "Commits"
    - Primary Key: CommitId
    - Columns: RepositoryId, Commit Date, Commit User, Repository Name
