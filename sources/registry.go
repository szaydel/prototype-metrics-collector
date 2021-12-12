package sources

type IMetricSource interface {
	Collect() (map[string]interface{}, error)
	Initialize(conf map[string]interface{}) error
	Name() string
}

type MetricSourceFactory func() IMetricSource

var Sources = map[string]MetricSourceFactory{}
