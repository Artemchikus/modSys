package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type Regression struct {
	data []*dataPoint
}

type dataPoint struct {
	Observed  float64
	Variables []float64
}

func (reg *Regression) Run() ([]float64, error) {
	numOfObservations := len(reg.data)
	numOfVars := len(reg.data[0].Variables)

	if numOfObservations < (numOfVars + 1) {
		return nil, fmt.Errorf("not enough observations: for %v variables (needs %v, got %v)", numOfVars+1, numOfVars+1, numOfObservations)
	}

	observed := mat.NewDense(numOfObservations, 1, nil)
	variables := mat.NewDense(numOfObservations, numOfVars+1, nil)

	for i := 0; i < numOfObservations; i++ {
		observed.Set(i, 0, reg.data[i].Observed)

		for j := 0; j < numOfVars+1; j++ {
			if j == 0 {
				variables.Set(i, 0, 1)
			} else {
				variables.Set(i, j, reg.data[i].Variables[j-1])
			}
		}
	}

	xTx := &mat.Dense{}
	xTx.Mul(variables.T(), variables)

	xTxI := &mat.Dense{}
	xTxI.Inverse(xTx)

	xTy := &mat.Dense{}
	xTy.Mul(variables.T(), observed)

	b := &mat.Dense{}
	b.Mul(xTxI, xTy)

	coffs := make([]float64, numOfVars+1)
	for i := 0; i < numOfVars+1; i++ {
		coffs[i] = b.At(i, 0)
	}

	return coffs, nil
}
