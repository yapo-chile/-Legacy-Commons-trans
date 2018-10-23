package infrastructure

import (
	"net/http"

	"github.com/newrelic/go-agent"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/loggers"
)

// NewRelicHandler struct representing a NewRelic handler with the new relic app
type NewRelicHandler struct {
	Appname string
	Key     string
	app     newrelic.Application
	Enabled bool
	Logger  loggers.Logger
}

// Start initializes the NewRelicHandler
func (n *NewRelicHandler) Start() error {
	if !n.Enabled {
		n.Logger.Info("NewRelic Off")
		return nil
	}
	config := newrelic.NewConfig(n.Appname, n.Key)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return err
	}
	n.app = app
	n.Logger.Info("NewRelic On")
	return err
}

// TrackHandlerFunc instruments an http.HandlerFunc
func (n *NewRelicHandler) TrackHandlerFunc(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	if !n.Enabled {
		return handler
	}
	return func(w http.ResponseWriter, r *http.Request) {
		txn := n.app.StartTransaction(pattern, w, r)
		defer txn.End() // nolint: errcheck
		handler(txn, r)
	}
}
