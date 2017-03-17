package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/skuid/helm-value-store-controller/controller"
	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/spec"
	"github.com/skuid/spec/lifecycle"
	_ "github.com/skuid/spec/metrics"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"time"
)

var portFlag = flag.IntP("port", "p", 3000, "The port to listen on")
var tableNameFlag = flag.StringP("table", "t", "helm-charts", "The DynamoDB table to read from")
var labelsFlag spec.SelectorSet
var intervalFlag = flag.StringP("interval", "i", "300s", "The sync interval to check the value store")
var blacklistFlag = flag.StringArrayP("blacklist", "b", []string{}, "A list of release names to not update or install")

func init() {
	var err error
	spec.Logger, err = spec.NewStandardLogger()
	if err != nil {
		fmt.Printf("Error setting up logger, %s\n", err.Error())
		os.Exit(1)
	}

}

func main() {
	flag.VarP(&labelsFlag, "labels", "l", "The labels to search the value store for.")
	flag.Parse()
	syncInterval, err := time.ParseDuration(*intervalFlag)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rs, err := dynamo.NewReleaseStore(*tableNameFlag)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go controller.SyncReleases(rs, labelsFlag.ToMap(), syncInterval, *blacklistFlag)

	muxer := http.NewServeMux()
	muxer.Handle("/metrics", promhttp.Handler())
	muxer.HandleFunc("/live", lifecycle.LivenessHandler)
	muxer.HandleFunc("/ready", lifecycle.ReadinessHandler)

	hostPort := fmt.Sprintf(":%d", *portFlag)
	spec.Logger.Info("Server is starting", zap.String("listen", hostPort))

	server := &http.Server{Addr: hostPort, Handler: muxer}
	lifecycle.ShutdownOnTerm(server)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		spec.Logger.Fatal("Error listening", zap.Error(err))
	}
	spec.Logger.Info("Server gracefully stopped")
}
