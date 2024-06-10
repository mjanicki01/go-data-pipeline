# Real-time Data Processing with Go, Elastic Beanstalk, and Redshift


## Project Overview

This project demonstrates a real-time data processing pipeline using a Go application. The application processes data stored in S3 and stores the results in Redshift. The Go application is deployed using Elastic Beanstalk.

## Objective

- Create a real-time data processing pipeline with a Go application.
- Deploy the application using Elastic Beanstalk.
- Store processed data in Amazon Redshift.
- Verify data processing and storage using SQL queries in Redshift.

## Prerequisites

Ensure you have the following installed on your local machine:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/dl/)

## .env Variables

```dosini
# To read from S3:
REGION=
BUCKET=
KEY=     #name of the .csv file in S3

# To push data to Redshift:
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
REDSHIFT_CONN_STRING=
```

## Run with Docker
```dosini
docker-compose build
docker-compose up

# To print the processed data:
curl "http://localhost:8080?action=print"
(note: the data is printed in the Docker container's console, not where curl is called)

# To insert processed data into Redshift:
curl "http://localhost:8080?action=insert"
```

## Run without Docker
```dosini
go mod tidy

# To print the processed data:
go run main.go -action=print

# To insert processed data into Redshift:
go run main.go -action=insert
```


