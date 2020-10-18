package main

import (
	"fmt"
	// "log"
	"net/http"
	"os"
	// "os/signal"
	// "runtime"
	// "syscall"
	// "text/tabwriter"
	"reflect"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "rds"
)

// Metrics descriptions
var (

	// labels are the static labels that come with every metric
	labels = []string{"region"}

	storage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "storage"),
		"Amount of storage for the RDS instance",
		labels, // labels
		nil,
	)
)

type promHTTPLogger struct {
	logger log.Logger
}

func (l promHTTPLogger) Println(v ...interface{}) {
	level.Error(l.logger).Log("msg", fmt.Sprint(v...))
}

type exporter struct {
	api    rdsiface.RDSAPI
}

// Describe describes the metrics exported by the RDS exporter. It
// implements prometheus.Collector.
func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- storage
}

// Collect fetches the stats from the configured RDS and delivers them
// as Prometheus metrics. It implements prometheus.Collector
func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		storage, prometheus.GaugeValue, 4.0, "us-east-1",
	)
}

func init() {
	prometheus.MustRegister(version.NewCollector("aws_rds_exporter"))
}

func run() int {

	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9107").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)


	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	fmt.Printf("Starting aws-rds-exporter...")
	fmt.Printf("\n")

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})

	if err != nil {
		fmt.Printf("Error getting session")
		return 1
	}

	if sess == nil {
		fmt.Printf("Error getting session after first error")
	} else {
		sType := reflect.TypeOf(sess)
		fmt.Printf(sType.String())
	}

	exporter := &exporter{api: rds.New(sess)}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath,
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer,
			promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					ErrorLog: &promHTTPLogger{
						logger: logger,
					},
				},
			),
		),
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>AWS RDS Exporter</title></head>
             <body>
             <h1>AWS RDS Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             <h2>Options</h2>
             <h2>Build</h2>
             </body>
             </html>`))
	})
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	fmt.Printf("\n")

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}

	return 0
}


func main() {

	exCode := run()
	os.Exit(exCode)
}




