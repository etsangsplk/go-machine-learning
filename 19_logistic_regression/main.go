package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kniren/gota/dataframe"

	"github.com/hfogelberg/golum"
)

const (
	scoreMax = 830.0
	scoreMin = 640.0
)

// Logistic regression from

// Goal: Determine if for a given credit score it is
// possible to get a loan with an interest rate below 12%.
func main() {
	if err := cleanData(); err != nil {
		return
	}

	df, err := golum.GetDFFromCSV("clean_loan_data.csv", nil)
	if err != nil {
		log.Printf("Error opening CSV file %s\n", err.Error())
		return
	}

	if err := createHistograms(&df); err != nil {
		return
	}

	if err := getStatistics(&df); err != nil {
		return
	}

	if err := split("clean_loan_data.csv"); err != nil {
		return
	}
}

func split(file string) error {
	_, _, err := golum.TrainTestSplit(file, 0.3)
	if err != nil {
		return err
	}

	return nil
}

func getStatistics(df *dataframe.DataFrame) error {
	cols := []string{"FICO.Range"}
	s, err := golum.GetStatistics(df, cols)
	if err != nil {
		return err
	}

	golum.PrintStats(s[0])

	return nil
}

func createHistograms(df *dataframe.DataFrame) error {
	// Create histogram for each column in the data set
	if err := golum.CreateHistograms(df, nil); err != nil {
		return err
	}

	return nil
}

func cleanData() error {
	f, err := os.Open("../data/loan_data.csv")
	if err != nil {
		log.Printf("Error opening file %s\n", err.Error())
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 2

	rawData, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading data %s\n", err.Error())
		return err
	}

	// Create output file
	f, err = os.Create("clean_loan_data.csv")
	if err != nil {
		log.Printf("Error creating output file %s\n", err.Error())
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)

	// Sequentially move the rows, writing out cleaned data
	for i, record := range rawData {
		// Write out header
		if i == 0 {
			if err := w.Write(record); err != nil {
				log.Printf("Error writing header %s\n", err.Error())
				return err
			}
			continue
		}

		// Initialize a slice to hold parsed values.
		outRecord := make([]string, 2)

		// Parse and standarize the FICO score
		score, err := strconv.ParseFloat(strings.Split(record[0], "-")[0], 64)
		if err != nil {
			log.Printf("Error parsing FICO score %s\n", err.Error())
			continue
		}
		outRecord[0] = strconv.FormatFloat((score-scoreMin)/(scoreMax-scoreMin), 'f', 4, 64)

		// Parse the interest rate class
		rate, err := strconv.ParseFloat(strings.TrimSuffix(record[1], "%"), 64)
		if err != nil {
			log.Printf("Error parsing interest rate %0.2f %s", rate, err.Error())
			continue
		}

		if rate <= 12.0 {
			outRecord[1] = "1.0"
		} else {
			outRecord[1] = "0.0"
		}

		if err := w.Write(outRecord); err != nil {
			fmt.Printf("Error writing to file %s\n", err.Error())
			return err
		}

		// Write buffered data
		w.Flush()

	}

	return nil
}
