package main

import (
	"fmt"

	"github.com/gonum/stat/distuv"

	"gonum.org/v1/gonum/stat"
)

func main() {
	observed := []float64{
		260.0,
		135.0,
		105.0,
	}

	totalObserved := 500.0

	// Calculate the expected frequencies
	expected := []float64{
		totalObserved * 0.60,
		totalObserved * 0.25,
		totalObserved * 0.15,
	}

	// Calculate ChiSquare to test statistics.
	// It will determine how well the data fits the expected results
	chiSquare := stat.ChiSquare(observed, expected)
	fmt.Printf("Chi-square is: %0.2f\n", chiSquare)

	// Calculate p-value
	// Create Chi square distribution with K degrees of freedom
	// K = 3-1=2 degrees of freedom
	// It is the number of possible categories minus 1.
	chiDist := distuv.ChiSquared{K: 2.0, Src: nil}
	pValue := chiDist.Prob(chiSquare)
	fmt.Printf("p-value: %0.4f\n", pValue)

	// P-value is 0.0001.
	// This means that there is a 0.01% chanse that the deviation is by pure chanse

}
