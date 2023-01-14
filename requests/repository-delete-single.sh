#/usr/bin/bash

curl -w 'Total: %{time_total}s\n' --location --request DELETE 'http://localhost:8080/api/glass/v1/repository/1'