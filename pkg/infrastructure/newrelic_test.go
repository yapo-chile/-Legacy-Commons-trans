package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewrelicStartError(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}

	nr := NewRelicHandler{
		Appname: "Test",
		Key:     "NotAValidKey",
		Enabled: true,
		Logger:  &mlogger,
	}
	err := nr.Start()
	assert.Error(t, err)
}

func MockHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("been there"))
}

func TestNewrelicStartOk(t *testing.T) {
	var conf NewRelicConf
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Info").Return(nil)
	LoadFromEnv(&conf)
	nr := NewRelicHandler{
		Appname: conf.Appname,
		Key:     conf.Key,
		Enabled: false,
		Logger:  &mlogger,
	}
	err := nr.Start()
	assert.NoError(t, err)

	m := MockHandlerFunc
	handler := nr.TrackHandlerFunc("test", m)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))
	handler(w, r)

	assert.Equal(t, "been there", w.Body.String())
	mlogger.AssertExpectations(t)
}
