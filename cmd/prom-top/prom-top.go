package main

import (
	"log"
	"fmt"
	"net/http"
	"time"
	"sync"
	"strconv"

	"github.com/gree-gorey/prom-top/pkg/cpu"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	chars [10]*prometheus.GaugeVec
	cpuUsage *prometheus.GaugeVec
)

func main() {
	for i, _ := range chars {
		chars[i] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: fmt.Sprintf("promtop_char%s", strconv.Itoa(i)),
				Help: fmt.Sprintf("char %s", strconv.Itoa(i)),
			},
			[]string{"rank"},
		)
		prometheus.MustRegister(chars[i])
	}

	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "promtop_cpu",
			Help: "cpu usage of a process",
		},
		[]string{"rank"},
	)
	prometheus.MustRegister(cpuUsage)

	http.Handle("/metrics", prometheus.Handler())
	go Run(int(3), true)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Run(interval int, debug bool) {
	c := cpu.CPU{Debug: false}
	for {
		var wg sync.WaitGroup
		wg.Add(1)
		go c.RunJob(&wg)
		wg.Wait()
		for _, proc := range c.CPUByProcess {
			cpuUsage.With(prometheus.Labels{"rank": strconv.Itoa(proc.Rank)}).Set(float64(proc.CPUUsedPercent))

			for i, _ := range proc.Chars {
				chars[i].With(prometheus.Labels{"rank": strconv.Itoa(proc.Rank)}).Set(float64(proc.Chars[i]))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
