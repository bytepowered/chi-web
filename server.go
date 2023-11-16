package chiweb

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	Server     *http.Server
	RootRouter chi.Router
}

func NewHttpServer(bindAddr string) *HttpServer {
	return &HttpServer{
		Server: &http.Server{
			Addr:         bindAddr,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		},
		RootRouter: newDefaultRootRouter(),
	}
}

func (s *HttpServer) Serve(serveCtx context.Context) error {
	s.Server.Handler = s.RootRouter
	log.Print("http server listen", "addr", s.Server.Addr)
	err := s.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HttpServer) Shutdown(shutdownCtx context.Context) error {
	slog.Info("http server shutdown")
	if err := s.Server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}

func (s *HttpServer) GracefulShutdown(serverCtx context.Context, timeout time.Duration) error {
	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, _ := context.WithTimeout(serverCtx, timeout)
	go func() {
		<-shutdownCtx.Done()
		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			log.Print("http server graceful shutdown timed out.. forcing exit.")
		}
	}()
	// Trigger graceful shutdown
	return s.Shutdown(shutdownCtx)
}

func GoServe(server *HttpServer, serverCtx context.Context, serverStopCtx context.CancelFunc) {
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		server.GracefulShutdown(serverCtx, time.Second*10)
		serverStopCtx()
	}()
	if err := server.Serve(serverCtx); err != nil {
		log.Print("web server serve", "error", err)
	}
	<-serverCtx.Done()
}

func SendJSON(w http.ResponseWriter, statusCode int, data []byte) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return w.Write(data)
}

func SendTEXT(w http.ResponseWriter, statusCode int, data []byte) (int, error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	return w.Write(data)
}

func newDefaultRootRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(10 * time.Second))
	return r
}
