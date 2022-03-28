// Schemes: https
// Host: watchdog.com
// BasePath: /api/v1
// Version: 1.0.0
// Contact: Habib MAALEM<habib.maalem@watchdog.com> https://www.watchdog.com
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Security:
// - bearer
// SecurityDefinitions:
// bearer:
//
//	type: apiKey
//	name: Authorization
//	in: header
//
// swagger:meta
// go:generate swagger generate spec
package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/internal/server/api/response"
	"github.com/groupe-edf/watchdog/internal/server/broadcast"
	"github.com/groupe-edf/watchdog/internal/server/job"
	"github.com/groupe-edf/watchdog/internal/server/middleware"
	"github.com/groupe-edf/watchdog/internal/server/queue"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/internal/version"
	"github.com/groupe-edf/watchdog/pkg/authentication"
	"github.com/groupe-edf/watchdog/pkg/authentication/token"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/logging"
	"github.com/pkg/errors"
)

const (
	apiRoot = "/api/v1"
)

var (
	ErrNotFound = errors.New("NOT_FOUND")
)

type API struct {
	logger     logging.Interface
	options    *config.Options
	store      models.Store
	workerPool *queue.WorkerPool
}

type RequestOptions struct {
	RequireAuthentication bool
}

func (api *API) Broadcast(scope string, message interface{}) {
	response := response.JSON(http.StatusOK, message)
	broadcast := container.GetContainer().Get(broadcast.ServiceName).(*broadcast.Broadcast)
	broadcast.Message <- response.Body()
}

func (api *API) NotFound(r *http.Request) response.Response {
	return response.Error(http.StatusNotFound, "", ErrNotFound)
}

func (api *API) Version(r *http.Request) response.Response {
	return response.JSON(http.StatusOK, version.GetBuildInfo())
}

// Mount the API's endpoints in the given router.
func (api *API) Mount(router *mux.Router) {
	wrap := func(next func(r *http.Request) response.Response, options RequestOptions) http.HandlerFunc {
		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := next(r)
			if result.Header() != nil {
				header := w.Header()
				for key, value := range result.Header() {
					header[key] = value
				}
			}
			api.SendResponse(w, result)
		})
		authentication := container.Get(authentication.ServiceName).(*authentication.Authentication)
		authenticator := middleware.NewAuthenticator(token.JWTOptions{
			BearerHeader:  authentication.Options.BearerHeader,
			Secret:        authentication.Options.Secret,
			JWTHeaderKey:  authentication.Options.JWTHeaderKey,
			JWTQueryParam: authentication.Options.JWTQueryParam,
		})
		authenticator.RequireAuthentication = options.RequireAuthentication
		return authenticator.Wrap(handlerFunc).ServeHTTP
	}
	protectedRoute := RequestOptions{
		RequireAuthentication: true,
	}
	router.HandleFunc("/-/health", wrap(api.Health, RequestOptions{})).Methods(http.MethodGet)
	router.HandleFunc("/-/oauth/redirect", wrap(api.OAuthLogin, RequestOptions{}))
	router.HandleFunc(apiRoot+"/", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/access_tokens", wrap(api.GetAccessTokens, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/access_tokens", wrap(api.SaveAccessToken, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/analytics", wrap(api.GetAnalytics, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/analyze", wrap(api.Analyze, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/analyzes", wrap(api.GetAnalyzes, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/analyzes/{analysis_id}", wrap(api.GetAnalysis, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/analyzes/{analysis_id}", wrap(api.DeleteAnalysis, protectedRoute)).Methods(http.MethodDelete)
	router.HandleFunc(apiRoot+"/authentication/forgot", wrap(api.Forgot, RequestOptions{})).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/authentication/login", wrap(api.Login, RequestOptions{})).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/authentication/register", wrap(api.Register, RequestOptions{})).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/authentication/reset", wrap(api.Reset, RequestOptions{})).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc(apiRoot+"/categories", wrap(api.GetCategories, RequestOptions{})).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/integrations", wrap(api.GetIntegrations, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/integrations", wrap(api.SaveIntegration, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/integrations/{integration_id}", wrap(api.GetIntegration, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/integrations/{integration_id}", wrap(api.DeleteIntegration, protectedRoute)).Methods(http.MethodDelete)
	router.HandleFunc(apiRoot+"/integrations/{integration_id}/groups", wrap(api.GetIntegrationGroup, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/integrations/{integration_id}/synchronize", wrap(api.SynchronizeIntegration, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/integrations/{integration_id}/webhooks", wrap(api.HandleWebhook, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/integrations/{integration_id}/webhooks", wrap(api.InstallWebhook, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/issues", wrap(api.GetIssues, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/jobs", wrap(api.GetJobs, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/jobs/{job_id}", wrap(api.GetJobs, protectedRoute)).Methods(http.MethodDelete)
	router.HandleFunc(apiRoot+"/leaks", wrap(api.GetLeaks, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/leaks/{leak_id}", wrap(api.GetLeak, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/pattern", wrap(api.EvaluatePattern, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/password", wrap(api.ChangePassword, protectedRoute)).Methods(http.MethodPut)
	router.HandleFunc(apiRoot+"/profile", wrap(api.GetProfile, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/policies", wrap(api.GetPolicies, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/policies", wrap(api.NewPolicy, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/policies/{policy_id}", wrap(api.GetPolicy, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/policies/{policy_id}", wrap(api.DeletePolicy, protectedRoute)).Methods(http.MethodDelete)
	router.HandleFunc(apiRoot+"/policies/{policy_id}/toggle", wrap(api.TogglePolicy, protectedRoute)).Methods(http.MethodPut)
	router.HandleFunc(apiRoot+"/policies/{policy_id}/conditions", wrap(api.AddPolicyCondition, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/policies/{policy_id}/conditions/{condition_id}", wrap(api.DeletePolicyCondition, protectedRoute)).Methods(http.MethodDelete)
	router.HandleFunc(apiRoot+"/queues", wrap(api.GetQueues, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/repositories", wrap(api.GetRepositories, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/repositories", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/repositories/{repository_id}", wrap(api.GetRepository, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/repositories/{repository_id}", wrap(api.DeleteRepository, protectedRoute)).Methods(http.MethodDelete)
	router.HandleFunc(apiRoot+"/repositories/{repository_id}/analyze", wrap(api.Analyze, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/repositories/{repository_id}/badge", wrap(api.GetRepositoryBadge, RequestOptions{})).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/rules", wrap(api.GetRules, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/rules", wrap(api.NewRule, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/rules/{rule_id}/toggle", wrap(api.ToggleRule, protectedRoute)).Methods(http.MethodPut)
	router.HandleFunc(apiRoot+"/settings", wrap(api.GetSettings, RequestOptions{})).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/settings", wrap(api.SaveSettings, protectedRoute)).Methods(http.MethodPost)
	router.HandleFunc(apiRoot+"/users", wrap(api.GetUsers, protectedRoute)).Methods(http.MethodGet)
	router.HandleFunc(apiRoot+"/version", wrap(api.Version, RequestOptions{})).Methods(http.MethodGet)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		broadcast := container.GetContainer().Get(broadcast.ServiceName).(*broadcast.Broadcast)
		broadcast.Upgrade(w, r)
	}).Methods(http.MethodGet, http.MethodOptions, http.MethodPost)
	router.NotFoundHandler = wrap(api.NotFound, RequestOptions{})
}

func (api *API) SendError(w http.ResponseWriter, err error) {
	problem := response.NewProblem()
	w.Header().Set("Content-Type", response.ProblemMediaType)
	switch errors.Cause(err).(type) {
	default:
		if err == sql.ErrNoRows || err == ErrNotFound {
			problem.Status = http.StatusNotFound
		} else {
			problem.Status = http.StatusInternalServerError
		}
		problem.Detail = err.Error()
	}
	if problem.ProblemStatus() != 0 {
		w.WriteHeader(problem.ProblemStatus())
	}
	response, _ := json.Marshal(problem)
	w.Write(response)
}

func (api *API) SendResponse(w http.ResponseWriter, data response.Response) {
	w.WriteHeader(data.Status())
	if data.Body() != nil {
		w.Write(data.Body())
	}
}

func NewAPI(ctx context.Context, options *config.Options, logger logging.Interface) (*API, error) {
	store, ok := container.GetContainer().Get(store.ServiceName).(models.Store)
	if !ok {
		return nil, errors.New("error when loading store")
	}
	// Setting up queue
	processor := job.ProcessAnalyze{
		Context: ctx,
		Logger:  logger,
		Options: options,
		Store:   store,
	}
	pool := queue.NewWorkerPool(
		queue.NewClient(store, logger),
		queue.WorkMap{
			"analyze_repository": processor.Handle,
		},
		queue.Options{
			MaxWorkers: 1,
		},
	)
	go pool.Start()
	return &API{
		logger:     logger,
		options:    options,
		store:      store,
		workerPool: pool,
	}, nil
}
