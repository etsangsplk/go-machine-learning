package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

func main() {
	// Create two matrices of the same size and a third with a different size
	a := mat.NewDense(3, 3, []float64{1, 2, 3, 0, 5, 6, 0, 0, 6})
	b := mat.NewDense(3, 3, []float64{8, 9, 10, 1, 4, 2, 9, 0, 2})
	c := mat.NewDense(3, 2, []float64{3, 2, 1, 4, 0, 8})

	// Add a and b
	d := mat.NewDense(0, 0, nil)
	d.Add(a, b)
	fd := mat.Formatted(d, mat.Prefix("            "))
	fmt.Printf("d = a + b = %0.4v\n\n", fd)

	// Multiply a and c
	f := mat.NewDense(0, 0, nil)
	f.Mul(a, c)
	ff := mat.Formatted(f, mat.Prefix("            "))
	fmt.Printf("f = a * c = %0.4v\n\n", ff)

	// Raise a matrix to a power
	g := mat.NewDense(0, 0, nil)
	g.Pow(a, 5)
	fg := mat.Formatted(g, mat.Prefix("            "))
	fmt.Printf("g = a ^ 5 = %0.4v\n\n", fg)

	// Apply a function to each elemt of a
	// You can apply any function to each element of an array
	h := mat.NewDense(0, 0, nil)
	sqrt := func(_, _ int, v float64) float64 { return math.Sqrt(v) }
	h.Apply(sqrt, a)
	fh := mat.Formatted(h, mat.Prefix("              "))
	fmt.Printf("h = sqrt(a) = %0.4v\n\n", fh)
}
