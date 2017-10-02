package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sajari/regression"

	"github.com/kniren/gota/dataframe"
)

func main() {
	f, err := os.Open("../data/Advertising.csv")
	if err != nil {
		log.Printf("Error opening file %s\n", err.Error())
		return
	}
	defer f.Close()

	df := dataframe.ReadCSV(f)

	_ = trainTestSplit(df)
	cols := []string{"TV", "Radio"}
	colsPos := []int{0, 1}
	if err := trainModel("Sales", cols, colsPos, 3); err != nil {
		log.Printf("Error training model %s\n", err.Error())
	}

}

func trainModel(yCol string, xCols []string, xPos []int, yPos int) error {
	f, err := os.Open("training.csv")
	if err != nil {
		fmt.Printf("Error opening CSV file %s\n", err.Error())
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4
	td, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file %s\n", err.Error())
		return err
	}

	var r regression.Regression
	r.SetObserved(yCol)

	for i, col := range xCols {
		r.SetVar(i, col)
	}

	// Loop CSV record and add training data
	for i, record := range td {
		if i == 0 {
			continue
		}

		yVal, err := strconv.ParseFloat(record[yPos], 64)
		if err != nil {
			log.Printf("Error parsing y values %s\n", err.Error())
			return err
		}

		var xVals []float64
		for _, pos := range xPos {
			fVal, err := strconv.ParseFloat(record[pos], 64)
			if err != nil {
				log.Printf("Error parsing x values %s\n", err.Error())
				return err
			}
			xVals = append(xVals, fVal)
		}

		r.Train(regression.DataPoint(yVal, xVals))
	}

	r.Run()
	fmt.Printf("Regression formula %s\n", r.Formula)

	return nil
}

func trainTestSplit(df dataframe.DataFrame) error {
	// Calculate number of elements in each set
	trainNum := (4 * df.Nrow()) / 5
	testNum := df.Nrow() / 5

	if trainNum+testNum < df.Nrow() {
		trainNum++
	}

	// Create subset indices
	trainIdx := make([]int, trainNum)
	testIdx := make([]int, testNum)

	// Enumerate indices
	for i := 0; i < trainNum; i++ {
		trainIdx[i] = i
	}

	for i := 0; i < testNum; i++ {
		testIdx[i] = i
	}

	traindDF := df.Subset(trainIdx)
	testDF := df.Subset(testIdx)

	// Create map that will be used in writing the data to files
	setMap := map[int]dataframe.DataFrame{
		0: traindDF,
		1: testDF,
	}

	// Create files
	for idx, setName := range []string{"training.csv", "test.csv"} {
		f, err := os.Create(setName)
		if err != nil {
			log.Printf("Error creating CSV file %s\n", err.Error())
			return err
		}

		w := bufio.NewWriter(f)
		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Printf("Error writing to CSV file %s\n", err.Error())
			return err
		}
	}

	log.Println("Train-test split done")
	return nil
}

func train() {

}
