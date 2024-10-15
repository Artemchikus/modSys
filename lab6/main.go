package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func hw(y []float64, alpha, beta, gamma float64, seasonLength, predictionLength int) ([]float64, error) {
	seasonals := make([]float64, seasonLength)
	seasonAverages := []float64{}
	nSeasons := len(y) / seasonLength

	for i := 0; i < nSeasons; i++ {
		seasonalAvg := 0.0
		for j := seasonLength * i; j < seasonLength*i+seasonLength; j++ {
			seasonalAvg += y[j] / float64(seasonLength)
		}

		seasonAverages = append(seasonAverages, seasonalAvg)
	}

	for i := 0; i < seasonLength; i++ {
		vAvg := 0.0

		for j := 0; j < nSeasons; j++ {
			vAvg += (y[seasonLength*j+i] - seasonAverages[j]) / float64(nSeasons)
		}

		seasonals[i] = vAvg
	}

	trend := 0.0
	for i := 0; i < seasonLength; i++ {
		trend += (y[i+seasonLength] - y[i]) / float64(seasonLength*seasonLength)
	}

	smooth := y[0]
	res := []float64{y[0]}

	for i := 1; i < len(y)+predictionLength; i++ {
		if i >= len(y) {
			m := float64(i - len(y) + 1)
			res = append(res, (smooth+m*trend)+seasonals[i%seasonLength])
		} else {
			lastSmooth := smooth
			smooth = alpha*(y[i]-seasonals[i%seasonLength]) + (1-alpha)*(smooth+trend)
			trend = beta*(smooth-lastSmooth) + (1-beta)*trend
			seasonals[i%seasonLength] = gamma*(y[i]-smooth) + (1-gamma)*seasonals[i%seasonLength]
			res = append(res, smooth+trend+seasonals[i%seasonLength])
		}
	}

	return res, nil
}

func main() {
	file, err := os.OpenFile("data.csv", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	transformer := unicode.BOMOverride(encoding.Nop.NewDecoder())

	reader := csv.NewReader(transform.NewReader(file, transformer))
	reader.Comma = ';'

	rec, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	y := make([]float64, len(rec[0]))
	for i, value := range rec[0] {
		flY, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}

		y[i] = flY
	}

	alpha := 0.02
	beta := 0.05
	gamma := 0.3
	seasonLength := 7
	m := 4

	info, err := hw(y, alpha, beta, gamma, seasonLength, m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("info: %v\n", info)
}
