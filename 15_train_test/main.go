package main

import (
	"bufio"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
)

// Do train-test split
// Todo: Create reusable package from the code?

func main() {
	f, err := os.Open("../data/diabetes.csv")
	if err != nil {
		log.Printf("Error opening file %s\n", err.Error())
		return
	}
	defer f.Close()

	df := dataframe.ReadCSV(f)
	trainNum := (4 * df.Nrow()) / 5
	testNum := df.Nrow() / 4
	if trainNum+testNum < df.Nrow() {
		trainNum++
	}

	// Create subset of indices
	trainIdx := make([]int, trainNum)
	testIdx := make([]int, testNum)

	// Enumerate training and testing indices
	for i := 0; i < trainNum; i++ {
		trainIdx[i] = i
	}

	for i := 0; i < testNum; i++ {
		testIdx[i] = i
	}

	// Create subset dataframes
	trainDF := df.Subset(testIdx)
	testDF := df.Subset(trainIdx)

	// Create map that will be used to write data files
	setMap := map[int]dataframe.DataFrame{
		0: trainDF,
		1: testDF,
	}

	// Create files
	for idx, setName := range []string{"training.csv", "test.csv"} {
		f, err := os.Create(setName)
		if err != nil {
			log.Printf("Error creating file %s\n", err.Error())
			return
		}
		w := bufio.NewWriter(f)
		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Printf("Error writing to file %s\n", err.Error())
			return
		}

	}

}
