package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yleong/pagerduty/pdapi"
)

func (s *Server) handleConfig(w *responseWithStatus, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(s)
	if err != nil {
		s.error(w, r, err)
		return
	}
	w.Write(b)
	return
}

func (s *Server) handleOncalls(w *responseWithStatus, r *http.Request) {
	schedules, err := s.PD.GetSchedules()
	if err != nil {
		s.error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(schedules.String()))
	return
}

func (s *Server) error(w *responseWithStatus, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	return
}

//Listen starts a HTTP server on the specified port
//The handlers return various pagerduty data
func (s *Server) Listen() {
	http.HandleFunc("/config", makeHandler(s.handleConfig))
	http.HandleFunc("/oncalls", makeHandler(s.handleOncalls))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", s.Port), nil))
}

//Server holds the HTTP and pagerduty configurations
type Server struct {
	Port string
	PD   pdapi.PagerDuty
}

var (
	clientRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_requests_total",
		Help: "The total number of HTTP requests received",
	}, []string{"path", "status_code"})
)
