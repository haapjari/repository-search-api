# Glass

This is **Glass**, a research tool which aims to offer functionality to measure quality of open-source repositories and libraries and return a single value called "Quality Measure" to represent state of repositories. Tool is used to create a data collection for my Master's thesis. Plugins represent the functionality, of what kind of repositories can be analyzed.

## Plugins

- *WIP*: Go

---

## How-To: Run

- See `Makefile`
- Requires: `go`, `postgresql`

---

## How-To: Contribute

- I don't have a structured way to accept contributions to the project, but feel free to leave a `pull request`, if you feel like it. :)

---

# Notes

## TODO

- Implement Grafana and Prometheus to collect stats on the research run.
- Create a Dockerfile.
- GitHub Actions Pipeline (?)
- Unit Tests.
- Secure the API.
- Deploy to the cloud. What cloud service provider will be used?
- Write some simple Swagger API documentation.
- Makefile -> Tasks (?)
- Create a release, when the first iteration is complete.

---

## Data Collection and Research

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

```
query {
        repository(owner: "TBD", name: "TBD") {
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
        }
}
```

- `script/list_commits.sh`
    - Amount of Collaborators in the Repository
    - Total Count of Commits in Repository
    - Every Commit in the Repository
    - Date of Commit in the Repository
    - Author of Commit in the Repository


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

#### Derivative Information

- Ratio: Library Codebase Size / Original Codebase Size
- Ratio: Open Issues / Closed Issues
- Commit Activity
- Maintentance: Collaborator Activity

---

### Database Tables

- Table: "Repository"
    - Primary Key: RepositoryId
    - Repository Struct: Repository Name, Url, CommitCount, Collaborators, Open Issues, Closed Issues, Original Codebase Size, Total Library Codebase Size, ProjectType, PrimaryLanguage
- Table: "Commits"
    - Primary Key: CommitId
    - Columns: RepositoryId, Commit Date, Commit User, Repository Name

---
