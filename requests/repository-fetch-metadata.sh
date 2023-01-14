#/usr/bin/bash

curl -w 'Total: %{time_total}s\n' --location --request GET 'http://localhost:8080/api/glass/v1/repository/fetch/metadata?type=go&count=1'