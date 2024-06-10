package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
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
		log.Printf("Successfully inserted %v into %s", key, tableName)
	}
	return nil
}

// handleAction performs the specified action (print or insert)
func handleAction(conn *pgx.Conn, config S3Config, action string) error {
	records, err := ReadCSVFromS3(config)
	if err != nil {
		return fmt.Errorf("Error reading CSV from S3: %v", err)
	}

	productSales, countrySales := ProcessData(records)

	switch action {
	case "print":
		PrintProcessedData(productSales, "Product Sales")
		PrintProcessedData(countrySales, "Country Sales")
	case "insert":
		err = InsertDataToRedshift(conn, "product_sales", productSales)
		if err != nil {
			return fmt.Errorf("Error inserting product sales data: %v", err)
		}
		err = InsertDataToRedshift(conn, "country_sales", countrySales)
		if err != nil {
			return fmt.Errorf("Error inserting country sales data: %v", err)
		}
	default:
		return fmt.Errorf("Invalid action specified. Please use -action=print or -action=insert")
	}
	return nil
}

func handleRequest(conn *pgx.Conn, config S3Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("action")
		err := handleAction(conn, config, action)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		fmt.Fprintf(w, "Action %s executed successfully", action)
	}
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

	action := flag.String("action", "", "Specify the action to perform: print or insert")
	flag.Parse()

	if *action == "" {
		http.HandleFunc("/", handleRequest(conn, config))
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		err := handleAction(conn, config, *action)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
