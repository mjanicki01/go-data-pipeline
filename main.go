package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type S3Config struct {
	Region string
	Bucket string
	Key    string
}

// Connect to S3 and get records from CSV file
func ReadCSVFromS3(config S3Config) ([][]string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region)},
	)
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(config.Key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// ProcessData calculates total sales for each product and each country
func ProcessData(records [][]string) (map[string]float64, map[string]float64) {
	productSales := make(map[string]float64)
	countrySales := make(map[string]float64)

	for _, record := range records[1:] {
		productID := record[1]
		country := record[7]
		quantity, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Println("Error parsing quantity:", err)
			continue
		}
		unitPrice, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			fmt.Println("Error parsing unit price:", err)
			continue
		}

		sales := quantity * unitPrice

		productSales[productID] += sales
		countrySales[country] += sales
	}

	return productSales, countrySales
}

// Print processed data into console (use for testing)
func PrintProcessedData(data map[string]float64, dataType string) {
	fmt.Printf("%s:\n", dataType)
	for key, value := range data {
		fmt.Printf("%s: %.2f\n", key, value)
	}
}

// Insert data into Redshift
func InsertDataToRedshift(conn *pgx.Conn, tableName string, data map[string]float64) error {
	for key, value := range data {
		_, err := conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO %s (id, total_sales) VALUES ($1, $2)", tableName), key, value)
		if err != nil {
			return fmt.Errorf("failed to insert data %v into %s: %v", key, tableName, err)
		}
		// log.Printf("Successfully inserted %v into %s", key, tableName)
	}
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connStr := os.Getenv("REDSHIFT_CONN_STRING")
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	log.Println("Successfully connected to the database")

	config := S3Config{
		Region: os.Getenv("REGION"),
		Bucket: os.Getenv("BUCKET"),
		Key:    os.Getenv("KEY"),
	}

	records, err := ReadCSVFromS3(config)
	if err != nil {
		fmt.Println("Error reading CSV from S3:", err)
		return
	}

	productSales, countrySales := ProcessData(records)

	// PrintProcessedData(productSales, "Product Sales")
	// PrintProcessedData(countrySales, "Country Sales")

	InsertDataToRedshift(conn, "product_sales", productSales)
	InsertDataToRedshift(conn, "country_sales", countrySales)

}
