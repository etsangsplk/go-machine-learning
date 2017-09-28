package main

/*
Example of parsing a CSV file and checking the data for missing values
without using third party packages.

If the CSV file is not delimited by commas and/or if it contains commented rows,
you can utilize the csv.Reader.Comma and csv.Reader.Comment fields to
*/

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type CSVRecord struct {
	SepalLength float64
	SepalWidth  float64
	PetalLength float64
	PetalWidth  float64
	Species     string
	ParseError  error
}

func main() {
	// Import and read CSV file
	f, err := os.Open("../data/iris.csv")
	if err != nil {
		log.Printf("Error opening iris data file %s\n", err.Error())
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// Each record should have five fields
	reader.FieldsPerRecord = 5
	var csvData []CSVRecord

	// Check for unexpected number of
	for {
		record, err := reader.Read()
		if err == io.EOF {
			log.Println("Done reading")
			break
		}

		var r CSVRecord
		// Parse each value in record based on expected type
		for idx, value := range record {
			if idx == 4 {
				// Skip if it's an empty string
				if value == "" {
					log.Printf("Unexpected type in column %d\n", idx)
					r.ParseError = fmt.Errorf("Empty string value")
					break
				}
				r.Species = value
				continue
			}

			// The rest of the values are float
			var fv float64

			// log and break if the value isn't a float
			if fv, err = strconv.ParseFloat(value, 64); err != nil {
				log.Printf("Unexpected type in column %s\n", err.Error())
				r.ParseError = fmt.Errorf("Could not parse float")
				break
			}

			// Add float value to correct column in CSV record
			switch idx {
			case 0:
				r.SepalLength = fv
			case 1:
				r.SepalWidth = fv
			case 2:
				r.PetalLength = fv
			case 3:
				r.PetalWidth = fv
			}

			if r.ParseError == nil {
				csvData = append(csvData, r)
			}
		}
	}

	for _, v := range csvData {
		fmt.Println(v)
	}
}
