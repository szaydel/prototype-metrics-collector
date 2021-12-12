package sources

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func makeName(words []string, index int) string {
	prefix := strings.Join(words[:1], ".")
	suffix := strings.Join(words[1:3], ".")
	return fmt.Sprintf("%s.%d.%s", prefix, index, suffix)
}

type MockSource struct {
	name string
}

func (s *MockSource) Collect() (map[string]interface{}, error) {
	log.Printf("%s Collect() method called", s.name)
	var mockData = map[string]interface{}{}

	for _, a := range []string{"alpha", "beta", "gamma", "delta", "epsilon"} {
		for _, b := range []string{"zeta", "eta", "theta", "iota", "kappa"} {
			for _, c := range []string{"lambda", "mu", "nu", "xi", "omicron"} {
				for i := 0; i < 5; i++ {
					mockData[makeName([]string{a, b, c}, i)] = rand.Float64()
				}
			}
		}
	}
	return mockData, nil
}

func (s *MockSource) Initialize(conf map[string]interface{}) error {
	return nil
}

func (s MockSource) Name() string {
	return s.name
}

func Register(name string, source MetricSourceFactory) {
	Sources[name] = source
}

func init() {
	s := &MockSource{name: "mocksource"}
	log.Printf("Registering %s\n", s.name)
	Register(s.name, func() IMetricSource { return s })
}
