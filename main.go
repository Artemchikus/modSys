package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"
)

func f(n int) []float64 {
	var b int32 = int32(time.Now().Unix())
	var two int32 = 2

	var y int32
	r := make([]float64, n)

	for i := 0; i < n; i++ {
		y = b * 1220703125

		if y < 0 {
			y += two * 1073741824
		}

		r[i] = float64(y) * 0.0000000004656613
		b = y
	}

	return r
}

func smirnoff(n []float64) float64 {
	slices.Sort(n)

	d := 0.0
	for i := 1; i <= len(n); i++ {
		tmpD := math.Abs(float64(i)/float64(len(n)) - n[i-1])
		if tmpD > d {
			d = tmpD
		}
	}

	return d
}

func chiSquare(n []float64, dof int) float64 {
	expected := float64(len(n)) / float64(dof)

	intervals := make([]float64, dof)

	for i := 0; i < len(n); i++ {
		for j := 1; j <= dof; j++ {
			if float64(j)/float64(dof) > n[i] {
				intervals[j-1]++
				break
			}
		}
	}

	chiSqr := 0.0
	for _, v := range intervals {
		chiSqr += (float64(v) - expected) * (float64(v) - expected) / expected
	}

	return chiSqr
}

var chiSquareMap = map[int]float64{
	1:   3.84,
	2:   5.99,
	3:   7.81,
	4:   9.49,
	5:   11.1,
	6:   12.6,
	7:   14.1,
	8:   15.5,
	9:   16.9,
	10:  18.3,
	11:  19.7,
	12:  21.0,
	13:  22.4,
	14:  23.7,
	15:  25.0,
	16:  26.3,
	17:  27.6,
	18:  28.9,
	19:  30.1,
	20:  31.4,
	21:  32.7,
	22:  33.9,
	23:  35.2,
	24:  36.4,
	25:  37.7,
	26:  38.9,
	27:  40.1,
	28:  41.3,
	29:  42.6,
	30:  43.8,
	40:  55.8,
	50:  67.5,
	60:  79.1,
	70:  90.5,
	80:  101.9,
	90:  113.1,
	100: 124.3,
}

var dofMap = map[int]float64{
	1:  0.975,
	2:  0.842,
	3:  0.708,
	4:  0.624,
	5:  0.565,
	6:  0.521,
	7:  0.486,
	8:  0.457,
	9:  0.432,
	10: 0.410,
	11: 0.391,
	12: 0.375,
	13: 0.361,
	14: 0.349,
	15: 0.338,
	16: 0.328,
	17: 0.318,
	18: 0.309,
	19: 0.301,
	20: 0.294,
	25: 0.27,
	30: 0.24,
	35: 0.23,
}

func chi(n int) (float64, error) {
	chi, ok := chiSquareMap[n-1]
	if ok {
		return chi, nil
	}

	return 0, fmt.Errorf("invalid degrees of freedom: %v", n)
}

func dof(n int) (float64, error) {
	dof, ok := dofMap[n]
	if ok {
		return dof, nil
	}

	if n < 35 {
		return 0, fmt.Errorf("invalid degrees of freedom: %v", n)
	}

	return 1.36 / math.Sqrt(float64(n)), nil
}

func main() {
	fl := f(18)
	fmt.Printf("fl: %v\n", fl)

	fl2 := make([]float64, len(fl))
	for i := 0; i < len(fl); i++ {
		fl2[i] = rand.Float64()
	}
	fmt.Printf("fl2: %v\n", fl2)

	dofA, _ := dof(18)
	fmt.Printf("dofA: %v\n", dofA)

	fmt.Printf("smirnov(fl): %v\n", smirnoff(fl))
	fmt.Printf("smirnov(fl2): %v\n", smirnoff(fl2))

	chSqr, _ := chi(10)
	fmt.Printf("chiSqr: %v\n", chSqr)
	chi1 := f(100)
	fmt.Printf("chi1: %v\n", chi1)

	chi2 := make([]float64, len(chi1))
	for i := 0; i < len(chi1); i++ {
		chi2[i] = rand.Float64()
	}
	fmt.Printf("chi2: %v\n", chi2)

	fmt.Printf("chiSquare(chi1, 100): %v\n", chiSquare(chi1, 10))
	fmt.Printf("chiSquare(chi2, 100): %v\n", chiSquare(chi2, 10))
}
