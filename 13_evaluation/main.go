package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/gonum/stat"
)

func main() {
	f, err := os.Open("../data/continuous_data.csv")
	if err != nil {
		log.Printf("Error opening data file %s\n", err.Error())
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)

	var observed []float64
	var predicted []float64
	i := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		// Skip the header
		if i == 1 {
			i++
			continue
		}

		obsVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("Parsing line %d failed %s\n", i, err.Error())
			continue
		}

		predVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Parsing line %d failed %s\n", i, err.Error())
			continue
		}

		observed = append(observed, obsVal)
		predicted = append(predicted, predVal)
		i++
	}

	// Calculate mean absolute error and mean square error
	var mAE float64
	var mSE float64

	for idx, val := range observed {
		mAE += math.Abs(val-predicted[idx]) / float64(len(observed))
		mSE += math.Pow(val-predicted[idx], 2) / float64(len(observed))
	}

	fmt.Printf("MAE: %0.2f\n", mAE)
	fmt.Printf("MSE: %0.2f\n", mSE)

	// MAE: 2.55
	// MSE: 10.51
	// The mean value of observed is 14.0, so MAE is just 20% of the observed

	// Calculate R-squared
	// Gives an idea about the deviation of the prediction
	// It measures the proportion of the variance in the observed values
	// in the predicted values

	rSquared := stat.RSquaredFrom(observed, predicted, nil)
	fmt.Printf("R-Squared: %0.2f\n", rSquared)
	// R-Squared: 0.37
	// Higher is better, so 0.37 is low
}
