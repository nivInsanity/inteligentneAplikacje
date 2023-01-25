package knna

import (
	"math"
	"sort"
)

type Knn struct {
	data []Data
	nrml []float64
}

func NewKnn(data []Data) *Knn {
	inputs := len(data[0].X)
	// normalization (for better comparison)
	nrml := make([]float64, inputs)
	mins := make([]float64, inputs)
	maxs := make([]float64, inputs)
	for i := 0; i < inputs; i++ {
		mins[i] = math.MaxFloat64
		maxs[i] = -math.MaxFloat64
	}
	for i := 0; i < len(data); i++ {
		x := data[i].X
		for j := 0; j < inputs; j++ {
			if x[j] < mins[j] {
				mins[j] = x[j]
			}
			if x[j] > maxs[j] {
				maxs[j] = x[j]
			}
		}
	}
	for i := 0; i < inputs; i++ {
		nrml[i] = maxs[i] - mins[i]
	}
	return &Knn{data, nrml}
}

func (knn *Knn) Distance(a []float64, b []float64) float64 {
	dist := 0.0
	for i := 0; i < len(a); i++ {
		dist += math.Abs(a[i]-b[i]) / knn.nrml[i]
	}
	return dist
}

func (knn *Knn) Predict(x []float64, k int) float64 {
	tosort := make([]Data, len(knn.data))
	copy(tosort, knn.data)
	sort.Slice(tosort, func(i, j int) bool {
		return knn.Distance(tosort[i].X, x) < knn.Distance(tosort[j].X, x)
	})
	predict_value := 0.0
	for i := 0; i < k; i++ {
		predict_value += tosort[i].Y
	}
	return predict_value / float64(k)
}
