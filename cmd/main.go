package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//docker exec -it app sh

//metrica de usuarios,gauge variacao conforme o tempo passa
var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_online_users",
	Help: "Online users",
	ConstLabels: map[string]string{//propiedade que se quer mostrar
		"logged_users": "true",
	},
})

//metrica para total de requests
var httpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "goapp_http_request_total",
	Help: "count all http requests for goapp",
}, []string{})

//metrica tempo de duracao do request, histogram
var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "goapp_http_requets_duration",
	Help: "http duration in seconds",
},[]string{"handler"})

func main() {
	//criando metricas para o prometheus
	r := prometheus.NewRegistry()
	r.MustRegister(onlineUsers)
	r.MustRegister(httpRequestTotal)
	r.MustRegister(httpDuration)

	go func ()  {
		for {
			//inserindo usuarios para simular metricas
			onlineUsers.Set(float64(rand.Intn(2000)))
		}
	}()

	home:= http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello welcome"))
	})
	//tempo de request na home
	duration := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "home"}),
		promhttp.InstrumentHandlerCounter(httpRequestTotal,home),
	)
	//verificar as o numero de requests em home
	http.Handle("/",duration)

	http.Handle("/metrics", promhttp.HandlerFor(r,promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8181",nil))
}