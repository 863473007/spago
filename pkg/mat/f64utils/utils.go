// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64utils

import (
	"github.com/nlpodyssey/spago/pkg/global"
	"gonum.org/v1/gonum/floats"
	"math"
	"strconv"
	"strings"
)

func EqualApprox(a, b float64) bool {
	return floats.EqualWithinAbsOrRel(a, b, 1.0e-06, 1.0e-06)
}

func Copy(in []float64) []float64 {
	out := make([]float64, len(in))
	copy(out, in)
	return out
}

// Fill the given slice's elements with value
func FillFloatSlice(slice []float64, value float64) {
	for i := range slice {
		slice[i] = value
	}
}

func Sign(a float64) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

func Max(v []float64) (m float64) {
	m = v[len(v)-1]
	for _, e := range v {
		if m <= e {
			m = e
		}
	}
	return
}

func Sum(v []float64) (s float64) {
	for _, e := range v {
		s += e
	}
	return
}

// ArgMinMax finds the indices of min and max arguments
func ArgMinMax(v []float64) (imin, imax int) {
	if len(v) < 1 {
		return
	}
	vmin, vmax := v[0], v[0]
	imin, imax = 0, 0
	for i := 1; i < len(v); i++ {
		if v[i] < vmin {
			imin = i
			vmin = v[i]
		}
		if v[i] > vmax {
			imax = i
			vmax = v[i]
		}
	}
	return
}

// ArgMax finds the index of the max argument
func ArgMax(v []float64) int {
	_, imax := ArgMinMax(v)
	return imax
}

// ArgMin finds the index of the min argument
func ArgMin(v []float64) int {
	imin, _ := ArgMinMax(v)
	return imin
}

func MakeFloat64Matrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, cols)
	}
	return matrix
}

func StrToFloat64Slice(str string) ([]float64, error) {
	spl := strings.Fields(str)
	data := make([]float64, len(spl))
	for i, v := range spl {
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			data[i] = num
		} else {
			return nil, err
		}
	}
	return data, nil
}

func SoftMax(v []float64) (sm []float64) {
	c := Max(v)
	var sum float64 = 0
	for _, e := range v {
		sum += math.Exp(e - c)
	}
	sm = make([]float64, len(v))
	for i, v := range v {
		sm[i] = math.Exp(v-c) / sum
	}
	return sm
}

// Tanh returns the hyperbolic tangent of x.
// If optimization level is 0, uses the default math.Tanh().
// If optimization level is 1, uses the "VarietyOfSound" tanh approximation.
// If optimization level is 2, uses the "Anguita" tanh approximation.
func Tanh(x float64) float64 {
	switch global.MathOptimizationLevel() {
	case 1:
		return tanhVarietyOfSound(x)
	case 2:
		return tanhAnguita(x)
	default: // no optimization
		return math.Tanh(x)
	}
}

// tanhAnguita calculate the tanh using the approximation described in:
// "Speed Improvement of the Back-Propagation on Current Generation Workstations" by Anguita et al, 1993.
func tanhAnguita(x float64) float64 {
	switch {
	case x > 1.92033:
		return 0.96016
	case x > 0:
		return 0.96016 - 0.26037*(x-1.92033)*(x-1.92033)
	case x <= -1.92033:
		return -0.96016
	default: // x < 0
		return 0.26037*(x+1.92033)*(x+1.92033) - 0.96016
	}
}

// tanhVarietyOfSound calculate the tanh using the approximation described in:
// https://varietyofsound.wordpress.com/2011/02/14/efficient-tanh-computation-using-lamberts-continued-fraction/
func tanhVarietyOfSound(x float64) float64 {
	x2 := x * x
	a := x * (135135.0 + x2*(17325.0+x2*(378.0+x2)))
	b := 135135.0 + x2*(62370.0+x2*(3150.0+x2*28.0))
	return a / b
}
