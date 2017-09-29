package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/kniren/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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
	describeData(df)

	if err := createHistograms(df); err != nil {
		return
	}

	if err := creatScatterplots(df); err != nil {
		return
	}

	if err := trainTestSplit(df); err != nil {
		return
	}

	if err := trainModel(); err != nil {
		return
	}

	if err := testModel(); err != nil {
		return
	}

	fmt.Println("Done!")
}

func testModel() error {
	f, err := os.Open("test.csv")
	if err != nil {
		log.Printf("Error opening test file %s\n", err.Error())
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

	// Loop over data and validate with Mean Absolute Error
	var mAE float64
	for i, record := range testData {
		// Skip header
		if i == 0 {
			continue
		}

		// Parse observed y (Sales)
		yObs, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("Error parsing y %s\n", err.Error())
			return err
		}

		// Parse X (TV value)
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("Error parsing X %s\n", err.Error())
			return err
		}

		// Predict Y
		var r regression.Regression
		yPred, err := r.Predict([]float64{tvVal})

		// Add to MAE
		mAE += math.Abs(yObs-yPred) / float64(len(testData))
	}
	fmt.Printf("MAE = %f0.2\n\n", mAE)

	return nil
}

func trainModel() error {
	f, err := os.Open("training.csv")
	if err != nil {
		fmt.Printf("Error opening training file %s\n", err.Error())
		return err
	}
	defer f.Close()

	// Read all training data
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading training data %s\n", err.Error())
		return err
	}

	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")
	for i, record := range trainingData {
		// Skip header
		if i == 0 {
			continue
		}

		// Parse sales val
		yVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("Error parsing sales data %s\n", err.Error())
			return err
		}

		// Parse TV val
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("Error parsing sales data %s\n", err.Error())
			return err
		}

		// Add data ponts to training model and do the training
		r.Train(regression.DataPoint(yVal, []float64{tvVal}))
		r.Run()
	}
	fmt.Printf("\nRegression formula: \n%v\n", r.Formula)
	return nil
}

func describeData(df dataframe.DataFrame) {
	// 	summary := df.Describe()
	// 	fmt.Println(summary)
}

// Do 80-20 split
func trainTestSplit(df dataframe.DataFrame) error {
	// Calculate number of elements in each set
	trainNum := (4 * df.Nrow()) / 5
	testNum := df.Nrow() / 5

	if trainNum+testNum < df.Nrow() {
		trainNum++
	}

	// Create subset indices
	trainIdx := make([]int, trainNum)
	testIdx := make([]int, testNum)

	// Enumerate indices
	for i := 0; i < trainNum; i++ {
		trainIdx[i] = i
	}

	for i := 0; i < testNum; i++ {
		testIdx[i] = i
	}

	traindDF := df.Subset(trainIdx)
	testDF := df.Subset(testIdx)

	// Create map that will be used in writing the data to files
	setMap := map[int]dataframe.DataFrame{
		0: traindDF,
		1: testDF,
	}

	// Create files
	for idx, setName := range []string{"training.csv", "test.csv"} {
		f, err := os.Create(setName)
		if err != nil {
			log.Printf("Error creating CSV file %s\n", err.Error())
			return err
		}

		w := bufio.NewWriter(f)
		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Printf("Error writing to CSV file %s\n", err.Error())
			return err
		}
	}

	log.Println("Train-test split done")
	return nil
}

// Scatterplot is a way of identifying the dependant variables
func creatScatterplots(df dataframe.DataFrame) error {
	// Ectract the target column
	yVals := df.Col("Sales").Float()
	for _, name := range df.Names() {
		// pts holds data to be plotted
		pts := make(plotter.XYs, df.Nrow())
		for i, fVal := range df.Col(name).Float() {
			pts[i].X = fVal
			pts[i].Y = yVals[i]
		}

		// Create plot
		p, err := plot.New()
		if err != nil {
			log.Printf("Error plotting %s %s\n", name, err.Error())
			return err
		}

		p.X.Label.Text = name
		p.Y.Label.Text = "y"
		p.Add(plotter.NewGrid())

		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Printf("Error creating new scatter plot for %s %s\n", name, err.Error())
			return err
		}

		s.GlyphStyle.Radius = vg.Points(3)
		p.Add(s)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, name+"_scatter.png"); err != nil {
			log.Printf("Error saving scatter plot %s\n", err.Error())
			return err
		}
	}

	log.Println("Done creating scatter plots")
	// Both Radio and TV have a somewhat linear relationship with sales
	return nil
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
			log.Printf("Error plotting %s %s\n", name, err.Error())
			return err
		}
		h.Normalize(1)

		// Add histogram to plot
		p.Add(h)
		if err := p.Save(4*vg.Inch, 4*vg.Inch, name+"_histogram.png"); err != nil {
			log.Printf("Error saving plot %s\n", err.Error())
			return err
		}
	}

	log.Println("Done creating histograms")
	return nil
}
