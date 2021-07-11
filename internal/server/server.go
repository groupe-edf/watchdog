package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/server/pool"
	v1 "github.com/groupe-edf/watchdog/internal/server/api/v1"
	"github.com/groupe-edf/watchdog/web"
	"golang.org/x/net/netutil"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	api      *v1.API
	context  context.Context
	pool     pool.Pool
	lock     sync.RWMutex
	logger   logging.Interface
	quitChan chan struct{}
	router   *mux.Router
}

func (server *Server) Close() {
	close(server.quitChan)
}

// Listener creates the TCP listener for web requests.
func (server *Server) Listener() (net.Listener, error) {
	listener, err := net.Listen("tcp", "0.0.0.0:3001")
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

func (server *Server) Run(ctx context.Context, listener net.Listener) error {
	server.context = ctx
	ctx, cancel := context.WithCancel(ctx)
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server.router)
	httpServer := &http.Server{
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		Handler: serveMux,
	}
	errChan := make(chan error, 1)
	go func() {
		server.logger.Info("starting server")
		errChan <- httpServer.Serve(listener)
	}()
	cancelChan := make(chan struct{})
	go func() {
		// Termination handler.
		termination := make(chan os.Signal, 1)
		signal.Notify(termination, os.Interrupt, syscall.SIGTERM)
		defer func() {
			signal.Stop(termination)
			cancel()
		}()
		select {
		case <-cancelChan:
			server.Close()
		case <-server.Quit():
			server.logger.Error("Received termination request via web service, exiting gracefully...")
			return
		case <-termination:
			server.logger.Error("Received SIGTERM, exiting gracefully...")
			server.Close()
			cancel()
			return
		}
	}()
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		httpServer.Shutdown(ctx)
		return nil
	case <-server.Quit():
		httpServer.Shutdown(ctx)
		return nil
	}
}

func New(logger logging.Interface) *Server {
	var clients = make(map[*websocket.Conn]bool)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		clients[connection] = true
	}).Methods("GET", "OPTIONS", "POST")
	server := &Server{
		pool: pool.Pool{
			Broadcast: make(chan []byte),
			Clients:   clients,
		},
		logger:   logger,
		quitChan: make(chan struct{}),
		router:   router,
	}
	server.api = v1.NewAPI(logger, server.pool)
	server.api.Register(server.router)
	webHandler := web.AssetsHandler("/", "build")
	server.router.Path("/").Handler(webHandler)
	server.router.Path("/manifest.json").Handler(webHandler)
	server.router.Path("/favicon.ico").Handler(webHandler)
	server.router.PathPrefix("/static/").Handler(webHandler)
	for _, route := range web.Routes {
		server.router.Path(route).Handler(webHandler)
	}
	go server.pool.Run()
	return server
}
