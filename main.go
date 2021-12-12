package main

import (
	// "log"
	"net/http"
	// "os"
	// "time"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/szaydel/prototype-metrics-collector/sources"
)

type App struct {
	collectors []sources.IMetricSource
	srvr       *echo.Echo
}

func NewServer(sources []sources.IMetricSource) *App {
	return &App{
		collectors: sources,
		srvr:       echo.New(),
	}
}

func (a *App) MetricsHandler() {
	a.srvr.GET("/metrics", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		var resultsChan = make(chan map[string]interface{})
		defer close(resultsChan)

		for _, s := range a.collectors {
			go func(c chan map[string]interface{}, source sources.IMetricSource) {
				data, err := source.Collect()
				if err != nil {
					a.srvr.Logger.Errorf("Failed to collect from '%s': %v", source.Name(), err)
					return
				}
				c <- data
			}(resultsChan, s)
		}

		var measurements = map[string]interface{}{}
		for i := 0; i < len(a.collectors); i++ {
			for key, value := range <-resultsChan {
				measurements[key] = value
			}
		}
		return json.NewEncoder(c.Response()).Encode(measurements)
	})
}

func (a *App) Start() {
	a.srvr.Logger.Fatal(a.srvr.Start(":8080"))
}

func main() {
	var collectors []sources.IMetricSource

	// Setup collectors, then start new HTTP server.
	for _, collector := range sources.Sources {
		s := collector()
		s.Initialize(map[string]interface{}{})
		collectors = append(collectors, s)
	}

	srvr := NewServer(collectors)
	srvr.MetricsHandler()
	srvr.Start()
}
