package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.schibsted.io/Yapo/goms/pkg/infrastructure"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/handlers"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/loggers"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/repository"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
)

var shutdownSequence = infrastructure.NewShutdownSequence()

func main() {
	var conf infrastructure.Config
	shutdownSequence.Listen()
	infrastructure.LoadFromEnv(&conf)
	if jconf, err := json.MarshalIndent(conf, "", "    "); err != nil {
		fmt.Printf("Config: \n%s\n", jconf)
	}

	fmt.Printf("Setting up logger\n")
	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf)
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
	// FibonacciHandler
	fibonacciLogger := loggers.MakeFibonacciInteractorLogger(logger)
	fibonacciRepository := repository.NewMapFibonacciRepository()
	fibonacciInteractor := usecases.FibonacciInteractor{
		Logger:     fibonacciLogger,
		Repository: fibonacciRepository,
	}
	fibonacciHandler := handlers.FibonacciHandler{
		Interactor: &fibonacciInteractor,
	}

	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger:        logger,
		WrapperFunc:   newrelic.TrackHandlerFunc,
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
						Name:    "Retrieve the Nth Fibonacci with Clean Architecture",
						Method:  "GET",
						Pattern: "/fibonacci",
						Handler: &fibonacciHandler,
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
	go server.ListenAndServe()
	shutdownSequence.Wait()

	logger.Info("Starting request serving")
	logger.Crit("%s\n", http.ListenAndServe(conf.ServiceConf.Host, maker.NewRouter()))
}
