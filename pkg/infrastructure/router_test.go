package infrastructure

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gorilla/mux.v1"
)

func TestRouterWithProfiling(t *testing.T) {
	var match mux.RouteMatch
	profiling := []bool{false, true}

	for _, with := range profiling {
		maker := RouterMaker{WithProfiling: with}
		router := maker.NewRouter()
		req := httptest.NewRequest("GET", "/debug/pprof/trace", strings.NewReader(""))
		doesMatch := router.Match(req, &match)
		assert.Equal(t, with, doesMatch)
	}
}
