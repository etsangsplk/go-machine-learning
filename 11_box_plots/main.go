package main

import (
	"fmt"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"github.com/kniren/gota/dataframe"
	"github.com/skratchdot/open-golang/open"
	// "gonum.org/v1/plot"
	// "gonum.org/v1/plot/plotter"
	// "gonum.org/v1/plot/vg"
)

func main() {
	iris, err := os.Open("../data/labeled_iris.csv")
	if err != nil {
		log.Printf("Error opening CSV file %s\n", err.Error())
		return
	}
	defer iris.Close()

	df := dataframe.ReadCSV(iris)
	fmt.Println(df)

	p, err := plot.New()
	if err != nil {
		log.Printf("Error creating plot %s\n", err.Error())
		return
	}

	p.Title.Text = "Box plots"
	p.Y.Label.Text = "Values"
	w := vg.Points(50)

	for idx, colName := range df.Names() {
		if colName != "species" {
			// Create plotter values and fill it with data
			v := make(plotter.Values, df.Nrow())
			for i, fv := range df.Col(colName).Float() {
				v[i] = fv
			}
			b, err := plotter.NewBoxPlot(w, float64(idx), v)
			if err != nil {
				log.Printf("Error plotting %s\n", err.Error())
				return
			}

			p.Add(b)
		}

		// Set x axis of plot to nominal with given names
		p.NominalX("sepal_length", "sepal_width", "petal_length", "petal_width")
		if err := p.Save(6*vg.Inch, 8*vg.Inch, "box_plots.png"); err != nil {
			log.Printf("Error saving file %s\n", err.Error())
			return
		}

		open.Run("box_plots.png")
	}

}
