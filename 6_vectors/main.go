package main

import (
	"fmt"

	"gonum.org/v1/gonum/floats"
)

// Use gonum.org/v1/gonum/floats or gonum.org/v1/gonum/math to handle vectors.

func main() {
	vectorA := []float64{11.0, 5.2, -1.3}
	vectorB := []float64{-7.2, 4.2, 5.1}

	dotProduct := floats.Dot(vectorA, vectorB)
	fmt.Printf("The dot product of A and B is %02f\n", dotProduct)

	// Scale vector A
	floats.Scale(1.5, vectorA)
	fmt.Printf("Scaling Vector A by 1.5: %v\n", vectorA)

	// Compute norm/length of vector B
	normB := floats.Norm(vectorB, 2)
	fmt.Printf("The norm/length of B is %0.2f\n", normB)
}
