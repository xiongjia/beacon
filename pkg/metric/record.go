package metric

type (
	DataPoint struct {
		Value     float64
		Timestamp int64
	}

	Record struct {
		Metric string
		Labels []Label
		DataPoint
	}
)
