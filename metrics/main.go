package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	http.Handle("/metrics", promhttp.Handler())

	// ...

	println("listening..")
	http.ListenAndServe(":5005", nil)
}
