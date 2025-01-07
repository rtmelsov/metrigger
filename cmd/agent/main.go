package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

type metrics map[string]float64

var ReportInterval int
var PollInterval int

func ParseFlag() {

	flag.IntVar(&ReportInterval, "r", 10, "report interval")
	flag.IntVar(&PollInterval, "p", 2, "poll interval")

	flag.Parse()
}

func main() {

	ParseFlag()

	met := make(chan metrics)

	go func(m chan metrics) {
		for {
			time.Sleep(time.Duration(PollInterval) * time.Second)
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			var met = metrics{}
			met["Alloc"] = float64(memStats.Alloc)
			met["BuckHashSys"] = float64(memStats.BuckHashSys)
			met["Frees"] = float64(memStats.Frees)
			met["GCCPUFraction"] = float64(memStats.GCCPUFraction)
			met["GCSys"] = float64(memStats.GCSys)
			met["HeapAlloc"] = float64(memStats.HeapAlloc)
			met["HeapIdle"] = float64(memStats.HeapIdle)
			met["HeapInuse"] = float64(memStats.HeapInuse)
			met["HeapObjects"] = float64(memStats.HeapObjects)
			met["HeapReleased"] = float64(memStats.HeapReleased)
			met["HeapSys"] = float64(memStats.HeapSys)
			met["LastGC"] = float64(memStats.LastGC)
			met["Lookups"] = float64(memStats.Lookups)
			met["MCacheInuse"] = float64(memStats.MCacheInuse)
			met["MCacheSys"] = float64(memStats.MCacheSys)
			met["MSpanInuse"] = float64(memStats.MSpanInuse)
			met["MSpanSys"] = float64(memStats.MSpanSys)
			met["Mallocs"] = float64(memStats.Mallocs)
			met["NextGC"] = float64(memStats.NextGC)
			met["NumForcedGC"] = float64(memStats.NumForcedGC)
			met["NumGC"] = float64(memStats.NumGC)
			met["OtherSys"] = float64(memStats.OtherSys)
			met["PauseTotalNs"] = float64(memStats.PauseTotalNs)
			met["StackInuse"] = float64(memStats.StackInuse)
			met["StackSys"] = float64(memStats.StackSys)
			met["Sys"] = float64(memStats.Sys)
			met["TotalAlloc"] = float64(memStats.TotalAlloc)
			m <- met
		}
	}(met)
	for {
		time.Sleep(time.Duration(ReportInterval) * time.Second)
		for k, b := range <-met {
			RequestToServer("counter", k, 1)
			RequestToServer("gauge", k, b)
		}
	}
}

func RequestToServer(t string, key string, value float64) {
	url := fmt.Sprintf("http://localhost:8080/update/%v/%v/%v", t, key, value)
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		log.Panic(err.Error())
	}

	req.Header.Add("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err.Error())
	}
	resp.Body.Close()
}
