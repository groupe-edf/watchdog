package server

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/core/models"
	v1 "github.com/groupe-edf/watchdog/internal/server/api/v1"
	"github.com/groupe-edf/watchdog/internal/server/broadcast"
	"github.com/groupe-edf/watchdog/internal/server/database"
	"github.com/groupe-edf/watchdog/internal/server/metrics"
	"github.com/groupe-edf/watchdog/internal/server/middleware"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/event"
	"github.com/groupe-edf/watchdog/pkg/logging"
	"github.com/groupe-edf/watchdog/web"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/netutil"
)

type Server struct {
	api        *v1.API
	context    context.Context
	cancelFunc context.CancelFunc
	di         container.Container
	logger     logging.Interface
	options    *config.Options
	quitChan   chan struct{}
	router     *mux.Router
	stopOnce   sync.Once
	server     *http.Server
}

func (server *Server) Close() {
	close(server.quitChan)
}

// Listener creates the TCP listener for web requests.
func (server *Server) Listener() (net.Listener, error) {
	listener, err := net.Listen("tcp", server.options.Server.ListenAddress)
	if err != nil {
		return listener, err
	}
	listener = netutil.LimitListener(listener, 32)
	return listener, nil
}

// Quit returns the receive-only quit channel.
func (server *Server) Quit() <-chan struct{} {
	return server.quitChan
}

func (server *Server) RegisterEvents() {
	bus := container.GetContainer().Get(event.ServiceName).(*event.EventBus)
	bus.SubscribeCallback("analysis:*", func(topic string, data interface{}) {
		analysis := data.(*models.Analysis)
		store := container.GetContainer().Get(store.ServiceName).(models.Store)
		_, _ = store.SaveAnalysis(analysis)
		broadcast := container.GetContainer().Get(broadcast.ServiceName).(*broadcast.Broadcast)
		broadcast.Broadcast(map[string]interface{}{
			"container_id":   analysis.ID.String(),
			"container_type": "analysis",
			"event_type":     topic,
			"payload":        analysis,
		})
	})
}

func (server *Server) SetupMetrics() {
	if server.options.Server.Profile {
		// Add the pprof routes
		server.router.HandleFunc("/debug/pprof/", pprof.Index)
		server.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		server.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		server.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		server.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
		server.router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
		server.router.Handle("/debug/pprof/block", pprof.Handler("block"))
		server.router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		server.router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		server.router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(database.QueryCounter)
	registry.MustRegister(database.QueryDuration)
	registry.MustRegister(metrics.HttpRequestHistogram)
	server.router.Handle("/-/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}

func (server *Server) Start(listener net.Listener) error {
	handler := http.NewServeMux()
	handler.Handle("/", server.router)
	middlewares := []middleware.Interface{
		&middleware.Logger{
			Logger: server.logger,
		},
		&middleware.CORS{},
		&metrics.Midleware{},
	}
	server.server = &http.Server{
		BaseContext: func(_ net.Listener) context.Context {
			return server.context
		},
		Handler: middleware.Merge(middlewares...).Wrap(handler),
	}
	errChan := make(chan error, 1)
	go func() {
		server.logger.Info("starting server")
		errChan <- server.server.Serve(listener)
	}()
	cancelChan := make(chan struct{})
	go func() {
		// Termination handler.
		termination := make(chan os.Signal, 1)
		signal.Notify(termination, os.Interrupt, syscall.SIGTERM)
		defer func() {
			signal.Stop(termination)
			server.cancelFunc()
		}()
		select {
		case <-cancelChan:
			server.Close()
		case <-server.Quit():
			server.logger.Error("received termination request via web service, exiting gracefully...")
			return
		case <-termination:
			server.logger.Error("received SIGTERM, exiting gracefully...")
			server.Close()
			server.cancelFunc()
			return
		}
	}()
	select {
	case err := <-errChan:
		return err
	case <-server.context.Done():
		server.Stop()
		return nil
	case <-server.Quit():
		server.Stop()
		return nil
	}
}

func (server *Server) Stop() {
	server.stopOnce.Do(func() {
		if err := server.server.Shutdown(server.context); err != nil {
			server.logger.Errorf("failed to shutdown server error: %s", err)
		}
	})
}

func New(ctx context.Context, logger logging.Interface) *Server {
	ctx, cancelFunc := context.WithCancel(ctx)
	di := container.GetContainer()
	options := di.Get(config.ServiceName).(*config.Options)
	router := mux.NewRouter().StrictSlash(true)
	server := &Server{
		context:    ctx,
		cancelFunc: cancelFunc,
		di:         di,
		logger:     logger,
		options:    options,
		quitChan:   make(chan struct{}),
		router:     router,
	}
	api, err := v1.NewAPI(ctx, options, logger)
	if err != nil {
		logger.Fatal(err)
	}
	server.api = api
	server.api.Mount(server.router)
	webHandler := web.AssetsHandler("/", "build")
	server.router.Path("/").Handler(webHandler)
	server.router.Path("/manifest.json").Handler(webHandler)
	server.router.Path("/favicon.ico").Handler(webHandler)
	server.SetupMetrics()
	server.router.PathPrefix("/static/").Handler(webHandler)
	for _, route := range options.Server.Static.Routes {
		server.router.Path(route).Handler(webHandler)
	}
	broadcast := di.Get(broadcast.ServiceName).(*broadcast.Broadcast)
	go broadcast.Run()
	return server
}
