package models

type Metric struct {
	ID                string
	Average           float64
	Median            float64
	Variance          float64
	StandardDeviation float64
	Mode              float64
}
