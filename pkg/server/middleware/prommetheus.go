package middleware

//import (
//	"github.com/bobgo0912/b0b-common/pkg/server"
//	"github.com/gorilla/mux"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promauto"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"net/http"
//)
//
//var (
//	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
//		Name: "myapp_http_duration_seconds",
//		Help: "Duration of HTTP requests.",
//	}, []string{"path"})
//)
//
//// PrometheusMiddleware implements mux.MiddlewareFunc.
//func PrometheusMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		route := mux.CurrentRoute(r)
//		path, _ := route.GetPathTemplate()
//		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
//		next.ServeHTTP(w, r)
//		timer.ObserveDuration()
//	})
//}
//
//func PromHttp(router *server.MuxRouter) {
//	router.Path("/metrics").Handler(promhttp.Handler())
//}
