package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gonum/floats"
	"github.com/kniren/gota/dataframe"
	"github.com/montanaflynn/stats"
	"gonum.org/v1/gonum/stat"
)

func main() {
	iris, err := os.Open("../data/labeled_iris.csv")
	if err != nil {
		log.Printf("Error opening CSV file %s\n", err.Error())
		return
	}
	defer iris.Close()

	df := dataframe.ReadCSV(iris)

	// Get float values of sepal_length and do some calculations
	sepalLength := df.Col("sepal_length").Float()
	mean := stat.Mean(sepalLength, nil)
	mode, modeCount := stat.Mode(sepalLength, nil)
	median, err := stats.Median(sepalLength)
	if err != nil {
		log.Printf("Error calculating median %s\n", err.Error())
	}
	min := floats.Min(sepalLength)
	max := floats.Max(sepalLength)
	rangeVal := max - min
	variance := stat.Variance(sepalLength, nil)
	stDev := stat.StdDev(sepalLength, nil)

	// Sort the values
	inds := make([]int, len(sepalLength))
	floats.Argsort(sepalLength, inds)

	// Get the Quantiles
	// 25% Quantile represents values where 25% of the are distributed below
	// the measure and 75% above.
	quant25 :=
		stat.Quantile(0.25, stat.Empirical, sepalLength, nil)
	quant50 :=
		stat.Quantile(0.50, stat.Empirical, sepalLength, nil)
	quant75 :=
		stat.Quantile(0.75, stat.Empirical, sepalLength, nil)

	fmt.Printf("Mean value: %0.2f\n", mean)
	fmt.Printf("Mode value: %0.2f\n", mode)
	fmt.Printf("Mode count: %0.2f\n", modeCount)
	fmt.Printf("Median: %0.2f\n", median)
	fmt.Printf("Min: %0.2f\n", min)
	fmt.Printf("Max: %0.2f\n", max)
	fmt.Printf("Range: %0.2f\n", rangeVal)
	fmt.Printf("Variance: %0.2f\n", variance)
	fmt.Printf("Standard deviation: %0.2f\n", stDev)
	fmt.Printf("25 Quantile: %0.2f\n", quant25)
	fmt.Printf("50 Quantile: %0.2f\n", quant50)
	fmt.Printf("75 Quantile: %0.2f\n", quant75)
}
