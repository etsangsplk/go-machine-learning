package main

import (
	"fmt"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"github.com/kniren/gota/dataframe"
)

func main() {
	f, err := os.Open("../data/Advertising.csv")
	if err != nil {
		log.Printf("Error opening file %s\n", err.Error())
		return
	}
	defer f.Close()

	df := dataframe.ReadCSV(f)

	// Describe not yet available
	describe(df)

	if err := createHistograms(df); err != nil {
		return
	}

}

Pfunc describe(df dataframe.DataFrame) {
	// 	summary := df.Describe()
	// 	fmt.Println(summary)
}

// Create histogram for each column in the data set
func createHistograms(df dataframe.DataFrame) error {
	for _, name := range df.Names() {
		plotVals := make(plotter.Values, df.Nrow())
		for i, fVal := range df.Col(name).Float() {
			plotVals[i] = fVal
		}

		// Make plot
		p, err := plot.New()
		if err != nil {
			log.Printf("Error creating plot %s\n", err.Error())
			return err
		}

		// Set title
		p.Title.Text = fmt.Sprintf("Histogram of a %s", name)

		// Create histogram
		h, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			fmt.Printf("Error plotting %s %s\n", name, err.Error())
			return err
		}
		h.Normalize(1)

		// Add histogram to plot
		p.Add(h)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, name+".png"); err != nil {
			fmt.Printf("Error saving plot %s\n", err.Error())
			return err
		}
	}

	fmt.Println("Done creating histograms")
	return nil
}
