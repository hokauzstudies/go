package prom

type metrictype struct {
	Type   string
	Name   string
	Helper string
	Labels []string
}

type counter struct{ metrictype }
type summary struct{ metrictype }
type histogram struct{ metrictype }

type collector struct {
	Indentifier string
	Metrics     []metrictype
}

var col = collector{
	Indentifier: "requests",
	Metrics: []metrictype{
		{"counter", "_endpoint_count", "The amount of request in this endpoint", []string{"code"}},
		{"summary", "_summary_durantion_seconds_", "The average of latency of HTTP of requests", []string{"handler", "method", "code"}},
		{"histogram", "_duration_deconds", "The latency of the HTTP requests", []string{"handler", "method", "code"}},
	},
}
