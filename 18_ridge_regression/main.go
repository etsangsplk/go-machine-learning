package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/berkmancenter/ridge"
	"github.com/gonum/matrix/mat64"
	"github.com/hfogelberg/golum"
)

func main() {
	trainFile, testFile, err := golum.TrainTestSplit("../data/Advertising.csv", 0.3)
	if err != nil {
		return
	}

	if err := train(trainFile); err != nil {
		return
	}

	if err := test(testFile); err != nil {
		return
	}

	fmt.Println("Done!")

}

func test(testFile string) error {
	f, err := os.Open(testFile)
	if err != nil {
		log.Printf("Error opening tes file %s\n", err.Error())
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4
	testData, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading test data %s\n", err.Error())
		return err
	}

	// Loop and calculate MAE
	var mAE float64
	for i, record := range testData {
		if i == 0 {
			continue
		}

		yObs, _ := strconv.ParseFloat(record[3], 64)
		tvVal, _ := strconv.ParseFloat(record[0], 64)
		radioVal, _ := strconv.ParseFloat(record[1], 64)
		paperVal, _ := strconv.ParseFloat(record[2], 64)

		// Predict y with trained model
		yPred := predict(tvVal, radioVal, paperVal)

		// Add to MAE
		mAE += math.Abs(yObs-yPred) / float64(len(testData))

	}

	fmt.Printf("\nMAE: %0.2f\n\n", mAE)

	return nil
}

func predict(tvVal float64, radioVal float64, paperVal float64) float64 {
	y := 3.297 + (0.043 * tvVal) + (0.202 * radioVal) + (-0.012 * paperVal)
	fmt.Printf("yPred: %0.3f\n", y)
	return y
}

func train(trainFile string) error {
	f, err := os.Open(trainFile)
	if err != nil {
		log.Printf("Error opening traing file %s\n", err.Error())
		return err
	}

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4

	rawData, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading csv %s\n", err.Error())
		return err
	}

	featureData := make([]float64, 4*len(rawData))
	yData := make([]float64, len(rawData))

	var fIdx int
	var yIdx int

	for i, record := range rawData {
		if i == 0 {
			continue
		}

		for i, val := range record {
			valParsed, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Printf("Error parsing %s %s", val, err.Error())
				continue
			}

			if i < 3 {
				// Add and intercept to the model
				if i == 0 {
					featureData[fIdx] = 1
					fIdx++
				}

				// Add floats to slice of feature floats
				featureData[fIdx] = valParsed
				fIdx++
			}

			if i == 3 {
				// Add value to slice of y floats
				yData[yIdx] = valParsed
				yIdx++
			}
		}
	}

	// For matrices that will be input to regression
	features := mat64.NewDense(len(rawData), 4, featureData)
	y := mat64.NewVector(len(rawData), yData)

	// Create new ridge regression
	// Penalty value is 1.0
	r := ridge.New(features, y, 1.0)
	r.Regress()

	c1 := r.Coefficients.At(0, 0)
	c2 := r.Coefficients.At(1, 0)
	c3 := r.Coefficients.At(2, 0)
	c4 := r.Coefficients.At(3, 0)

	fmt.Printf("\nRegression formula: \n")
	fmt.Printf("y = %0.3f + %0.3f TV + %0.3f Radio + %0.3f Newspaper\n", c1, c2, c3, c4)

	// Prints out
	// y = 3.297 + 0.043 TV + 0.202 Radio + -0.012 Newspaper
	// Newspaper adds have no effect

	return nil
}
