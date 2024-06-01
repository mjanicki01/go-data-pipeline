package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

type S3Config struct {
	Region string
	Bucket string
	Key    string
}

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

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

	fmt.Println("Product Sales:")
	for product, sales := range productSales {
		fmt.Printf("%s: %.2f\n", product, sales)
	}

	fmt.Println("Country Sales:")
	for country, sales := range countrySales {
		fmt.Printf("%s: %.2f\n", country, sales)
	}
}
