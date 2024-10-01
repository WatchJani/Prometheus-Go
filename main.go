package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	reqCounter prometheus.Counter
	// reqLatency prometheus.Histogram
}

func NewMetrics() *PromMetrics {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "actor_msg_counter",
		Help: "count number of req",
	})

	prometheus.MustRegister(counter)

	return &PromMetrics{
		reqCounter: counter,
	}
}

type User struct {
	*PromMetrics
}

func New(prom *PromMetrics) *User {
	return &User{
		PromMetrics: prom,
	}
}

func (u *User) Inc(w http.ResponseWriter, r *http.Request) {
	u.reqCounter.Inc()
}

func main() {
	user := New(NewMetrics())
	http.HandleFunc("/", user.Inc)
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		for {
			_, err := http.Get("http://localhost:5000/")
			if err != nil {
				log.Println(err)
			}

			time.Sleep(time.Second)
		}
	}()

	http.ListenAndServe(":5000", nil)
}
