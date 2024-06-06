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

## Run with Docker

1. `docker-compose build`
2. `docker-compose up`

## Run without Docker

1. `go mod tidy`
2. `go run main.go`
