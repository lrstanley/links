// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const promNamespace = "links"

var (
	metricShortened = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "shortened_total",
		Help:      "Total amount of shortened links",
	})
	metricRedirected = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "redirected_total",
		Help:      "Total amount of redirections",
	})
)

func init() {
	// Register metrics here.
	prometheus.MustRegister(
		metricShortened,
		metricRedirected,
	)
}

func initMetrics(ctx context.Context, wg *sync.WaitGroup, errors chan<- error, r *chi.Mux) {
	if !conf.Prometheus.Enabled {
		return
	}

	// Bind to existing http address router.
	if conf.Prometheus.Addr == "" {
		debug.Printf("binding metrics endpoint to '%s'", conf.Prometheus.Endpoint)
		debug.Println("WARNING: binding prometheus to existing address/port means it's much easier to accidentally expose to the public!")
		r.Handle(conf.Prometheus.Endpoint, promhttp.Handler())
		return
	}

	// Custom http server specifically for metrics (makes it easier to firewall off).
	mux := chi.NewRouter()
	if conf.Proxy {
		mux.Use(middleware.RealIP)
	}

	mux.Use(middleware.Compress(5))
	mux.Use(middleware.DefaultLogger)
	mux.Use(middleware.Timeout(10 * time.Second))
	mux.Use(middleware.Throttle(5))
	mux.Use(middleware.Recoverer)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<a href="%s">see metrics</a>`, conf.Prometheus.Endpoint)
	})

	mux.Handle(conf.Prometheus.Endpoint, promhttp.Handler())

	srv := &http.Server{
		Addr:         conf.Prometheus.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		wg.Add(1)
		defer wg.Done()

		debug.Printf("initializing metrics server on %s", conf.Prometheus.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errors <- fmt.Errorf("metrics error: %v", err)
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		<-ctx.Done()

		debug.Printf("requesting metrics server to shutdown")
		if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			errors <- fmt.Errorf("unable to shutdown metrics server: %v", err)
		}
	}()
}
