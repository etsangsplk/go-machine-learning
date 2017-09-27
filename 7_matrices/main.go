package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func main() {
	// Create flat representation of matrix
	data := []float64{1.2, -5.7, -2.3, 7.5, 3.4, 1.9}

	// Form matrix
	mx := mat.NewDense(3, 2, data)

	// Format and print
	fa := mat.Formatted(mx, mat.Prefix(""))
	fmt.Printf("%v\n", fa)

	// Access value at a certain position
	val := mx.At(2, 1)
	fmt.Printf("The value at 2,1 is %v\n", val)

	// Get values in a column
	col := mat.Col(nil, 0, mx)
	fmt.Printf("The values in the first column: %v\n", col)

	// Get values in a row
	row := mat.Row(nil, 1, mx)
	fmt.Printf("The values in the second row are %v\n", row)

	// Modify a value
	mx.Set(0, 0, 11.2)
	fmt.Printf("Modified values \n%v\n", mx)
}
