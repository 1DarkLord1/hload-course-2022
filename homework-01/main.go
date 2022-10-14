package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dword"
	password = "admin"
	dbname   = "postgres"
)

var (
	createOpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_processed_ops_total",
		Help: "The total number of processed /create requests",
	})
	
	urlOpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "url_processed_ops_total",
		Help: "The total number of processed /<tinyurl> requests",
	})
	
	createOpTime = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "create_processed_op_time_mcs",
		Help: "The duration of /create request",
	})
	
	urlOpTime = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "url_processed_op_time_mcs",
		Help: "The duration of /<tinyurl> request",
	})
)

func setupRouter(conn *sql.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.PUT("/create", func(c *gin.Context) {
		start := time.Now()
		createPutHandler(conn, c)
		elapsed := time.Since(start).Microseconds()

		createOpsProcessed.Inc()
		createOpTime.Observe(float64(elapsed))
	})

	r.GET("/:url", func(c *gin.Context) {
		start := time.Now()
		urlGetHandler(conn, c)
		elapsed := time.Since(start).Microseconds()

		urlOpsProcessed.Inc()
		urlOpTime.Observe(float64(elapsed))
	})

	return r
}

func main() {
	fmt.Println(sql.Drivers())

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Println("Failed to open", err)
		panic("exit")
	}

	err = conn.Ping()

	if err != nil {
		fmt.Println("Failed to ping database", err)
		panic("exit")
	}

	err = createUrlStorageTable(conn)

	if err != nil {
		fmt.Println("Failed create table", err)
		panic("exit")
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe(":8088", nil)
	}()
	
	r := setupRouter(conn)
	r.Run(":8080")
}
