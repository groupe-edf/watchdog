package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/server/pool"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/internal/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrNotFound = errors.New("NOT_FOUND")
)

type AnalyzeRepositoryCommand struct {
	From          string `json:"from"`
	RepositoryURL string `json:"repository_url"`
	Since         string `json:"since,omitempty"`
	Until         string `json:"until,omitempty"`
}

type API struct {
	logger logging.Interface
	pool   pool.Pool
	store  *store.Store
}

type Response struct {
	Data  interface{} `json:"data"`
	Error error       `json:"error,omitempty"`
	Scope string      `json:"scope,omitempty"`
}

func (api *API) Broadcast(scope string, message interface{}) {
	response, err := json.Marshal(Response{
		Data:  message,
		Scope: scope,
	})
	if err != nil {
		api.logger.Error("error marshaling message")
	}
	api.pool.Broadcast <- []byte(response)
}

func (api *API) NotFound(r *http.Request) Response {
	return Response{
		Error: ErrNotFound,
	}
}

func (api *API) Version(r *http.Request) Response {
	return Response{
		Data: version.GetBuildInfo(),
	}
}

// Register the API's endpoints in the given router.
func (api *API) Register(router *mux.Router) {
	wrap := func(callback func(r *http.Request) Response) http.HandlerFunc {
		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Accept, Origin, Authorization")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			api.logger.WithFields(logging.Fields{
				"method":  r.Method,
				"request": r.RequestURI,
			}).Info()
			result := callback(r)
			if result.Error != nil {
				api.SendError(w, result.Error)
				return
			}
			if result.Data != nil {
				api.SendResponse(w, result.Data)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
		return handlerFunc
	}
	router.HandleFunc("/-/metrics", promhttp.Handler().ServeHTTP)
	router.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/api/v1/analyze", wrap(api.Analyze)).Methods("OPTIONS", "POST")
	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {})
	router.HandleFunc("/api/v1/issues", wrap(api.GetIssues))
	router.HandleFunc("/api/v1/repositories", wrap(api.GetRepositories)).Methods("GET")
	router.HandleFunc("/api/v1/repositories", func(w http.ResponseWriter, r *http.Request) {}).Methods("OPTIONS", "POST")
	router.HandleFunc("/api/v1/repositories/{repository_id}", wrap(api.GetRepository))
	router.HandleFunc("/api/v1/rules", wrap(api.GetRules))
	router.HandleFunc("/api/v1/version", wrap(api.Version))
	router.NotFoundHandler = wrap(api.NotFound)
}

func (api *API) SendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	switch err {
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}

func (api *API) SendResponse(w http.ResponseWriter, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		api.logger.Error("error writing response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func NewAPI(logger logging.Interface, connectionPool pool.Pool) *API {
	store, err := store.NewStore(logger, &config.Database{
		Host:     "localhost",
		Name:     "watchdog",
		Password: "watchdog",
		Port:     5432,
		Username: "watchdog",
	})
	if err != nil {
		logger.Fatalf("error creating store : %s", err)
	}
	return &API{
		logger: logger,
		pool:   connectionPool,
		store:  store,
	}
}
