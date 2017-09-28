package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
	"github.com/skratchdot/open-golang/open"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Gonum plot is on Mercurial, so it must be installed to install plotting

func main() {
	iris, err := os.Open("../data/labeled_iris.csv")
	if err != nil {
		log.Printf("Error opening CSV file %s\n", err.Error())
		return
	}
	defer iris.Close()

	df := dataframe.ReadCSV(iris)
	// fmt.Println(df)

	// Create histogram for each feature column in the dataset
	for _, colName := range df.Names() {
		if colName != "species" {
			v := make(plotter.Values, df.Nrow())
			for i, val := range df.Col(colName).Float() {
				v[i] = val
			}

			p, err := plot.New()
			if err != nil {
				log.Printf("Error making plot %s\n", err.Error())
				return
			}
			p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

			// Create histogram of drawn values from the standard normal
			h, err := plotter.NewHist(v, 16)
			if err != nil {
				log.Printf("Error creating histogram %s\n", err.Error())
				return
			}

			// Normalize the diagram.
			// This makes it possible to compare different distributions side by side
			h.Normalize(1)

			// Add histogram to plot
			p.Add(h)

			// Save plot to png file
			name := colName + ".png"
			if err := p.Save(4*vg.Inch, 4*vg.Inch, name); err != nil {
				log.Printf("Error saving file %s\n", err.Error())
				return
			}
			log.Printf("Open %s\n", name)
			open.Run(name)

		}
	}
}
