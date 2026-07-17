package server

import (
	"context"
	"erp/internal/controllers"
	"erp/internal/middlewares"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
}

var _ Router = (*Server)(nil)

func NewServer() (server Server) {
	server.mux = http.NewServeMux()
	return
}

func (Server) __internal() {}

func (server *Server) Use(middlewares ...middlewares.Middleware) {
	server.middlewares = append(server.middlewares, middlewares...)
}

func (server *Server) Add(controllers ...controllers.Controller) {
	for _, controller := range controllers {
		server.mux.Handle(controller.Pattern(), controller)
	}
}

func (server *Server) Group(groupFunc func(g *Group)) {
	group := NewGroup(server.mux)
	groupFunc(&group)
	group.Chain()
}
func (server *Server) Serve(addr string) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: server.chain(),
	}

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return httpServer.Shutdown(ctx)
	}
}

func (server *Server) chain() http.Handler {
	handler := http.Handler(server.mux)
	for i := len(server.middlewares) - 1; i >= 0; i-- {
		handler = server.middlewares[i](handler)
	}
	return handler
}
