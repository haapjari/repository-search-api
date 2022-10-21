# Glass

This is Glass, a REST API written in Go, which aims to offers functionality to research GitHub -repositories. This is intended to be a research tool, which is part of my Master's Thesis.

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

- PostgreSQL -> Redis, with Snapshots
- README.md -file.
- GitHub Actions Pipeline.
- Unit Tests.
- Create a Dockerfile.
- Secure the API.
- Optimize the API.
- Deploy to the cloud. What cloud service provider will be used?
- Write some simple Swagger API documentation.
- Implement a Logger (?)
- Implement Grafana and Prometheus to collect stats on the research run.
- Makefile -> Tasks (?)
- Create a release, when the first iteration is complete.
- Test the endpoints.
    - Repository
        - GET / - works
        - POST / - works
        - GET /:id - not working
        - PATCH /:id - not tested
        - DELETE /:id - not tested
    - Commit
        - GET / - not tested
        - POST / - not tested
        - GET /:id - not tested
        - PATCH /:id - not tested
        - DELETE /:id - not tested

---

## Data Collection and Research

### GraphQL

- SourceGraph API, Query: `lang:go select:repo repohasfile:go.mod count:100000`
- GitHub API

---

### Quality Measure

- Quality Measure:
    - Repository Activity (higher is better)
        - Amount of commits, with dates. (TODO, How-To Query with GraphQL)
    - Maintenance (higher is better)
        - Amount of collabolators (TODO, How-To Query with GraphQL)
    - Code Smells (https://github.com/securego/gosec) (less is better) (TODO, Create Algorithm)
        - Amount of code smells in the repository.
        - Code smells severity average.
    - Ratio of Open Issues to Closed Issues (less is better)
        - Amount of Open Issues
        - Amount of Closed Issues
- Thresholds of these amounts will be calculated, threshholds will be in values 0 - 5, where 2.5 is at middle of the amounts.
- These values will be averaged in a single QM value. Correlation will be calculated ratio of library to original code lines, or ratio of sizes. Is there a correlation between bigger ratio and quality measure.

### Usage of gosec

- RUN: `$ gosec -fmt=json -out=results.json -stdout ./...`
- Calculate the lengh of issues array inside the JSON.
- Save that amount to the database.

---

### Database Tables

- Table: "Repository"
    - Primary Key: RepositoryId
    - Repository Struct: Repository Name, Url, CommitCount, Collaborators, Open Issues, Closed Issues, Original Codebase Size, Total Library Codebase Size, ProjectType, PrimaryLanguage
- Table: "Commits"
    - Primary Key: CommitId
    - Columns: RepositoryId, Commit Date, Commit User, Repository Name

### Derivative Information

- Ratio: Library Codebase Size / Original Codebase Size
- Ratio: Open Issues / Closed Issues
- Commit Activity
- Maintentance: Collaborator Activity

---