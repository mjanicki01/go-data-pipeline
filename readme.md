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

## AWS Setup

1. **Download Dataset**
   - Download the Online Retail dataset from [UCI Machine Learning Repository](https://archive.ics.uci.edu/dataset/352/online+retail).

2. **Create an S3 Bucket**
   - Follow the instructions to [create an S3 bucket](https://docs.aws.amazon.com/AmazonS3/latest/user-guide/create-bucket.html).

3. **Upload the Dataset to S3**
   - Upload the `Online Retail` CSV file to the S3 bucket you created. Instructions can be found [here](https://docs.aws.amazon.com/AmazonS3/latest/user-guide/upload-objects.html).

4. **Create a Redshift Cluster**
   - Follow the steps to [create a Redshift cluster](https://docs.aws.amazon.com/redshift/latest/gsg/getting-started.html).

5. **Create a New Database within the Redshift Cluster**
   - Once your Redshift cluster is created, use the [AWS Management Console](https://docs.aws.amazon.com/redshift/latest/mgmt/working-with-snapshots-console.html) or [AWS CLI](https://docs.aws.amazon.com/redshift/latest/mgmt/managing-clusters-cli.html) to create a new database within the cluster.


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


