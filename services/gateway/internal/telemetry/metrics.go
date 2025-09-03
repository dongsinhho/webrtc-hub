package telemetry

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MetricsHandler() http.Handler { return promhttp.Handler() }
func Window() time.Duration        { return 1 * time.Minute }
