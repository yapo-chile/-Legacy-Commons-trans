package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.mpi-internal.com/Yapo/trans/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/trans/pkg/interfaces/handlers"
	"github.mpi-internal.com/Yapo/trans/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/trans/pkg/interfaces/repository/services"
	"github.mpi-internal.com/Yapo/trans/pkg/usecases"
)

var shutdownSequence = infrastructure.NewShutdownSequence()

func main() {
	var conf infrastructure.Config
	shutdownSequence.Listen()
	infrastructure.LoadFromEnv(&conf)
	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		fmt.Printf("Config: \n%s\n", jconf)
	}

	fmt.Printf("Setting up Prometheus\n")
	prometheus := infrastructure.MakePrometheusExporter(
		conf.PrometheusConf.Port,
		conf.PrometheusConf.Enabled,
	)

	shutdownSequence.Push(prometheus)
	fmt.Printf("Setting up logger\n")
	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf,
		prometheus.NewEventsCollector(
			"trans_service_events_total",
			"events tracker counter for trans service",
		),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	logger.Info("Setting up New Relic")
	newrelic := infrastructure.NewRelicHandler{
		Appname: conf.NewRelicConf.Appname,
		Key:     conf.NewRelicConf.Key,
		Enabled: conf.NewRelicConf.Enabled,
		Logger:  logger,
	}
	err = newrelic.Start()
	if err != nil {
		logger.Error("Error loading New Relic: %+v", err)
		os.Exit(2)
	}

	logger.Info("Initializing resources")

	// HealthHandler
	var healthHandler handlers.HealthHandler

	// transHandler
	transFactory := infrastructure.NewTextProtocolTransFactory(conf.Trans, logger)
	transRepository := services.NewTransRepo(transFactory)
	transLogger := loggers.MakeTransInteractorLogger(logger)
	transInteractor := usecases.TransInteractor{
		Repository: transRepository,
		Logger:     transLogger,
	}

	transHandler := handlers.TransHandler{
		Interactor: transInteractor,
	}
	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger: logger,
		WrapperFuncs: []infrastructure.WrapperFunc{
			newrelic.TrackHandlerFunc,
			prometheus.TrackHandlerFunc,
		},
		WithProfiling: conf.ServiceConf.Profiling,
		Routes: infrastructure.Routes{
			{
				// This is the base path, all routes will start with this prefix
				Prefix: "/api/v{version:[1-9][0-9]*}",
				Groups: []infrastructure.Route{
					{
						Name:    "Check service health",
						Method:  "GET",
						Pattern: "/healthcheck",
						Handler: &healthHandler,
					},
					{
						Name:    "Execute a trans request",
						Method:  "POST",
						Pattern: "/execute/{command}",
						Handler: &transHandler,
					},
				},
			},
		},
	}
	server := infrastructure.NewHTTPServer(
		fmt.Sprintf("%s:%d", conf.Runtime.Host, conf.Runtime.Port),
		maker.NewRouter(),
		logger,
	)
	shutdownSequence.Push(server)
	logger.Info("Starting request serving")
	go server.ListenAndServe()
	shutdownSequence.Wait()

	logger.Info("Server exited normally")
}
