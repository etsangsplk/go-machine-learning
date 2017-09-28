package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
)

// Example of handling CSV data with data frames.

func main() {
	iris, err := os.Open("../data/labeled_iris.csv")
	if err != nil {
		log.Printf("Error opening CSV file %s\n", err.Error())
		return
	}
	defer iris.Close()

	df := dataframe.ReadCSV(iris)
	fmt.Println(df)

	// Filter data frame to only see Iris versicolor
	filter := dataframe.F{
		Colname:    "species",
		Comparator: "==",
		Comparando: "Iris-versicolor",
	}

	vsDF := df.Filter(filter)

	if vsDF.Err != nil {
		log.Printf("Error filtering data %s\n", err.Error())
		return
	}

	// Only select sepal width and species column
	vsDF = df.Filter(filter).Select([]string{"sepal_width", "species"})
	fmt.Println(vsDF)
}
