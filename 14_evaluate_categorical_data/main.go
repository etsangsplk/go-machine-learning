package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Accuracy: Percentage of predictions that were right
// (TP+TN)/(TP+TN+FP+FN)
// Precision: Percentage of correct positive predictions
// TP/(TP+FP)
// Recall: Percentage of positve predictions that were identified as positive
// TP/(TP+FN)

func main() {
	f, err := os.Open("../data/labeled.csv")
	if err != nil {
		log.Printf("Error opening file %s\n", err.Error())
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)

	var observed []int
	var predicted []int
	i := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		// Skip headers
		if i == 1 {
			i++
			continue
		}

		obsVal, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Error parsing record %d %s\n", i, err.Error())
			continue
		}

		predVal, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("Error parsing record %d %s\n", i, err.Error())
			continue
		}

		observed = append(observed, obsVal)
		predicted = append(predicted, predVal)
		i++
	}

	// fmt.Println(observed)
	// fmt.Println(predicted)

	// Hold count of TP and TN values
	var truePosNeg int

	// Accumulate true results
	for idx, val := range observed {
		if val == predicted[idx] {
			truePosNeg++
		}
	}

	// Calculate the accuracy
	acc := float64(truePosNeg) / float64(len(observed))
	fmt.Printf("Accuracy: %0.2f\n", acc)

}
