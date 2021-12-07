package sources

type IMetricSource interface {
	Collect() (map[string]interface{}, error)
}

type MetricSourceFactory func() IMetricSource

var Inputs = map[string]MetricSourceFactory{}
