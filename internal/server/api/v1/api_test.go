package v1

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/container"
	"github.com/groupe-edf/watchdog/internal/server/query"
	"github.com/groupe-edf/watchdog/internal/server/store"
)

var (
	api    API
	router *mux.Router
)

func TestMain(m *testing.M) {
	setUpAll()
	exitStatus := m.Run()
	tearDownAll()
	os.Exit(exitStatus)
}

func basicAuth(email, password string) string {
	credentials := email + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(credentials))
}

func performRequest(router http.Handler, method, path string, payload []byte) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(method, path, bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func setUpAll() {
	di := container.GetContainer()
	options := &config.Options{
		Server: &config.Server{
			Database: &config.Database{
				Driver: "bolt",
				Bolt: config.BoltOptions{
					Path: "/tmp",
				},
			},
		},
	}
	di.Set(config.ServiceName, func(c container.Container) container.Service {
		return options
	})
	di.Provide(&logging.ServiceProvider{})
	di.Provide(&store.ServiceProvider{
		Options: options.Database,
	})
	logger := di.Get(logging.ServiceName).(logging.Interface)
	store := di.Get(store.ServiceName).(models.Store)
	api = API{
		logger:  logger,
		options: options,
		store:   store,
	}
	router = mux.NewRouter().StrictSlash(true)
	api.Mount(router)
}

func tearDownAll() {
	api.store.DeleteUsers(&query.Query{})
}
