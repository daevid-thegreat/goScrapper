package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"time"
)

func main() {
	fName := "data.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("error closing file: %v", err)
		}
	}(file)

	// Create a new writer.
	writer := csv.NewWriter(file)
	// Write any buffered data to the underlying writer (standard output).
	defer writer.Flush()

	// Write CSV header
	c := colly.NewCollector(
		colly.AllowedDomains("jumia.com.ng"),
	)
	c.UserAgent = "goScrapper/1.0 (https://daevidthegreat.com/)"

	// Limit the number of threads started by colly to two
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "jumia.com.ng",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Handle errors during scraping
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s\n", r.Request.URL)
		log.Printf("Request Headers: %+v\n", r.Request.Headers)
		log.Printf("Response Status Code: %d\n", r.StatusCode)
		log.Printf("Error: %v\n", err)
	})

	c.OnHTML(".c-prd", func(e *colly.HTMLElement) {
		productName := e.ChildText(".name")
		productPrice := e.ChildText(".curr")
		productOldPrice := e.ChildText(".old")
		fmt.Println(productName, productPrice, productOldPrice)
		err := writer.Write([]string{
			productName,
			productPrice,
			productOldPrice,
		})
		if err != nil {
			log.Fatalf("error writing record to csv: %v", err)
		}
	})

	err = c.Visit("https://www.jumia.com.ng/mlp-bluetti-store/")
	if err != nil {
		log.Fatalf("error visiting URL: %v", err)
	}

	log.Printf("Scraping finished, check file %q for results\n", fName)
	log.Print(c)
}
