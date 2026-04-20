package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ListenAndServe starts the health HTTP server on the given addr (e.g. ":9090").
// It blocks until ctx is cancelled, then performs a graceful shutdown.
func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	httpServer := &http.Server{
		Addr:        addr,
		Handler:     s.Handler(),
		ReadTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("healthcheck: listen: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutCtx)
	}
}
