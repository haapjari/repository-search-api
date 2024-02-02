# Greetings

Hello! I am Jari Haapasaari ([mail](mailto:haapjari@gmail.com)), and this is a research tool for my thesis, and my first Go project! 

This is a complete open-source project, which started from my interest to learn `go`. Feel free to `fork`, leave comment, mail me, or copy parts of the code to your own usage. This `README.md` -file contains notes for the project, tasks, and descriptive information of the logic behind the functionality.

---

## About

This repository contains replication package for the analysis, which was completed in my thesis, during Autumn 2023. 

Work consists of three different components and a docker-compose file that can be used to orchestrate the setup.

Project is called `glass`, because I's like magnifying glass to GitHub repositories. Tool is an abstraction for `GitHub REST API`, which collects data of open-source repositories and writes that data to database.

### Components

### UI

- Visualizes the results.
- Offers orchestration methods and ability to visualize the results.

### API

- REST API, that fetches data from GitHub.
- API also calculates the missing fields from the data.

### Visual

- Component reads data from the database, and  

### Database

- PostgreSQL database.

---