package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

func dispersion(values []float64) float64 {
	sumXi := 0.0
	sumXiSqr := 0.0

	for _, val := range values {
		sumXi += val
		sumXiSqr += val * val
	}

	n := float64(len(values))

	return sumXiSqr/n - (sumXi/n)*(sumXi/n)
}

func corelation(x, y []float64) float64 {
	xAvg := 0.0
	yAvg := 0.0
	for i := 0; i < len(x); i++ {
		xAvg += x[i] / float64(len(x))
		yAvg += y[i] / float64(len(y))
	}

	corelTop := 0.0
	corelBotX := 0.0
	corelBotY := 0.0
	for i := 0; i < len(x); i++ {
		corelTop += (x[i] - xAvg) * (y[i] - yAvg)
		corelBotX += (x[i] - xAvg) * (x[i] - xAvg)
		corelBotY += (y[i] - yAvg) * (y[i] - yAvg)
	}

	return corelTop / (math.Sqrt(corelBotX * corelBotY))
}

func rmse(x, y []float64) float64 {
	sum := 0.0
	for i := 0; i < len(x); i++ {
		sum += (x[i] - y[i]) * (x[i] - y[i])
	}

	return math.Sqrt(sum / float64(len(x)))
}

func rSq(x, y []float64) float64 {
	avgY := 0.0

	for _, val := range x {
		avgY += val / float64(len(x))
	}

	rTop := 0.0
	rBot := 0.0

	for i := 0; i < len(x); i++ {
		rTop += (x[i] - y[i]) * (x[i] - y[i])
		rBot += (x[i] - avgY) * (x[i] - avgY)
	}

	return 1 - rTop/rBot
}

func main() {
	dataPoints, valss, err := parseCsv("data.csv")
	if err != nil {
		panic(err)
	}

	r := Regression{
		data: dataPoints,
	}

	coffs, err := r.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("y = %v", coffs[0])
	for i, val := range coffs[1:] {
		fmt.Printf(" + (x%d * %v)", i+1, val)
	}
	fmt.Printf("\n")

	yDisp := dispersion(valss[0])

	for i, vals := range valss[1:] {
		disp := dispersion(vals)
		fmt.Printf("dispersion of x%v: %v\n", i+1, disp)

		if disp/yDisp > 1.56 {
			fmt.Printf("f-test for x%v: %v (failed)\n", i+1, disp/yDisp)
			continue
		}

		fmt.Printf("f-test for x%v: %v (passed)\n", i+1, disp/yDisp)
	}

	for i, vals := range valss[1:] {
		corel := corelation(vals, valss[0])
		if math.Abs(corel) < 0.9 {
			fmt.Printf("corelation of x%v and y: %v (failed)\n", i+1, corel)
			continue
		}
		fmt.Printf("corelation of x%v and y: %v (passed)\n", i+1, corel)
	}

	for i, vals := range valss[1:] {
		for j := 1; j < len(valss); j++ {
			if i+1 == j {
				continue
			}
			corel := corelation(vals, valss[j])
			if math.Abs(corel) < 0.9 {
				fmt.Printf("corelation of x%v and x%v: %v (failed)\n", i+1, j, corel)
				continue
			}
			fmt.Printf("corelation of x%v and x%v: %v (passed)\n", i+1, j, corel)
		}
	}

	modeledYs := make([]float64, len(dataPoints))
	for i := 0; i < len(dataPoints); i++ {
		modeledYs[i] = coffs[0]
	}

	for i, vals := range valss[1:] {
		for j, val := range vals {
			modeledYs[j] += val * coffs[i+1]
		}
	}

	rmse := rmse(valss[0], modeledYs)
	fmt.Printf("rmse: %v\n", rmse)

	rSq := rSq(valss[0], modeledYs)
	fmt.Printf("r^2: %v\n", rSq)
}

func parseCsv(path string) ([]*dataPoint, [][]float64, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'

	rec, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var dataPoints []*dataPoint
	for _, observed := range rec[1][1:] {
		obs, err := strconv.ParseFloat(observed, 64)
		if err != nil {
			continue
		}

		dataPoint := &dataPoint{
			Observed: obs,
		}

		dataPoints = append(dataPoints, dataPoint)
	}

	for _, values := range rec[2:] {
		for i, val := range values[1:] {
			flVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			dataPoints[i].Variables = append(dataPoints[i].Variables, flVal)
		}
	}

	values := make([][]float64, 0, len(rec)-1)

	for _, vs := range rec[1:] {
		vals := make([]float64, 0, len(vs)-1)
		for _, v := range vs[1:] {
			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			vals = append(vals, val)
		}
		values = append(values, vals)
	}

	return dataPoints, values, nil
}
