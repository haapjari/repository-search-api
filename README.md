# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and this is a research tool for my thesis, and my first Go project! 

This is a complete open-source project, which started from my interest to learn `go`. Feel free to `fork`, leave comment, mail me, or copy parts of the code to your own usage. This `README.md` -file contains notes for the project, tasks, and descriptive information of the logic behind the functionality.

---

## About

I am calling this project `glass`, because I's like magnifying glass to GitHub repositories. Tool is an abstraction for `GitHub REST API`, which collects data of open-source repositories and writes that data to database.

---

## Plugins

- `Glass` is designed to be modular, `pkg/plugins` folder represents what kind of repositories can be analyzed. I am working (at the moment of writing, 27.10.2022), on `goplg`, which aims to offer functionality to analyze the quality of repositories, which primary language is `go`. Alpha Version of the Plugin is Completed in 15.1.2023.

- Go (Alpha Version Released)

---

# How-To: Run

- See `Makefile`
- Requires: `go`, `postgresql`
- Example: `.env` -file.

```
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_HOST=
POSTGRES_PORT=
GITHUB_USERNAME=
GITHUB_TOKEN=
```

---