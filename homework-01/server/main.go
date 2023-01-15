package main

import (
	"net/http"
	"main/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	createOpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_processed_ops_total",
		Help: "The total number of processed PUT /create requests",
	})
	
	getOpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_processed_ops_total",
		Help: "The total number of processed GET /<tinyurl> requests",
	})
	
	createOpTime = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "create_processed_op_time_mcs",
		Help: "The duration of /create request",
	})
	
	getOpTime = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "url_processed_op_time_mcs",
		Help: "The duration of GET /<tinyurl> request",
	})
)

func setupRouter(service *Service) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.PUT("/create", func(c *gin.Context) {
		utils.MeasureTime(func() { service.createPutHandler(c) }, createOpsProcessed, createOpTime)
	})

	r.GET("/:url", func(c *gin.Context) {
		utils.MeasureTime(func() { service.urlGetHandler(c) }, createOpsProcessed, createOpTime)
	})

	return r
}

func main() {
	service := &Service{}
	service.init()

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe(":8088", nil)
	}()
	
	r := setupRouter(service)
	r.Run(":8080")
}
