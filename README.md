# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and I originally developed the initial version of this project as a part of my thesis - and as my first project written in `Go`. Later I cleared and cleaned up the project, as the original version also looked like my very first `Go` project. 

- Original Version Released: `20th October 2023`.
- Cleaner Version Released: `24th April 2024`.

If you're interested into reproduce the research, please see: [repository-analysis-orchestration](https://github.com/haapjari/repository-analysis-orchestration) repository.

***

## About

This tool is meant to be an abstraction for a set of different GitHub API's that offer ability to query metadata from different repositories. You can see the API's within this file: [repository_service.go](https://github.com/haapjari/repository-search-api/blob/main/internal/pkg/service/repository_service.go). For the OpenAPI, please refer to the [openapi.yaml](https://github.com/haapjari/repository-search-api/blob/main/docs/openapi.yaml).

**NOTE**: Third-Party LOC Reporting works only with projects, written in Go.

***

## How-To

### Run

- You require `go` and `make` to run this project. Tested with `go-1.22.0`.
- Setup `PORT` as Environment Variable, and execute `make run` or just `PORT=8080 make run`.

### Build and Run as a Docker Container

- Build the Image: `docker build -t repository-search-api:latest .`
- Run the Image (On the Host, for Simplicity): `docker run -idt -e PORT=8080 --network=host repository-search-api:latest`

### Example Query

- You need to have [Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) exposed as a `GITHUB_TOKEN` environment variable to be able to run this command.

```bash
curl "localhost:8080/api/v1/repos/search?firstCreationDate=2008-01-01&lastCreationDate=2009-01-01&language=Go&minStars=100&maxStars=1000&order=desc" --header "Authorization: Bearer $GITHUB_TOKEN"
```
