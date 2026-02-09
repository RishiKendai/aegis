package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1. Compute API requests counter
var ComputeRequestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "aegis_compute_requests_total",
		Help: "Total number of requests to /compute API",
	},
)

// 2. Preprocess API requests counter
var PreprocessRequestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "aegis_preprocess_requests_total",
		Help: "Total number of requests to /api/v1/preprocess API",
	},
)

// 3. Plagiarism computation duration histogram
var PlagiarismComputationDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "aegis_plagiarism_computation_duration_seconds",
		Help:    "Time taken to compute plagiarism for each request",
		Buckets: prometheus.DefBuckets,
	},
)

// 4. Invalid submissions counter with reason label
var InvalidSubmissionsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "aegis_invalid_submissions_total",
		Help: "Total number of invalid submissions",
	},
	[]string{"reason"}, // All current reasons: invalid_request_body, missing_drive_id, no_artifacts, astra_preprocess_error, mongo_no_candidate_reports, mongo_no_document_plagiarism_reports, failed_to_update_candidate_result
)

// 5. High plagiarisms detected counter (counts individual candidates per compute request)
var HighPlagiarismsDetected = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "aegis_high_plagiarisms_detected_total",
		Help: "Total number of high plagiarisms detected (candidates with RiskHighlySuspicious or RiskNearCopy) per /compute request",
	},
	[]string{"drive_id"}, // Label to track per driveID
)

// InitPrometheus initializes and registers all Prometheus metrics
func InitPrometheus() {
	prometheus.MustRegister(ComputeRequestsTotal)
	prometheus.MustRegister(PreprocessRequestsTotal)
	prometheus.MustRegister(PlagiarismComputationDuration)
	prometheus.MustRegister(InvalidSubmissionsTotal)
	prometheus.MustRegister(HighPlagiarismsDetected)
}

// MetricsHandler returns the Prometheus metrics HTTP handler
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
