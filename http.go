package mtls

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListenAndServe starts plain HTTP or mTLS HTTPS depending on MTLS_ENABLED.
func ListenAndServe(handler http.Handler, addr string) error {
	cfg := LoadFromEnv()
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	if !cfg.Enabled {
		return srv.ListenAndServe()
	}

	tlsCfg, err := ServerTLSConfig(cfg)
	if err != nil {
		return err
	}
	srv.TLSConfig = tlsCfg
	return srv.ListenAndServeTLS(cfg.ServerCertFile, cfg.ServerKeyFile)
}

// RunGin is a drop-in for gin.Engine.Run with optional mTLS.
func RunGin(engine *gin.Engine, addr string) error {
	if err := ListenAndServe(engine, addr); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// ConfigureServer applies mTLS to an existing http.Server when enabled.
func ConfigureServer(srv *http.Server) error {
	cfg := LoadFromEnv()
	if !cfg.Enabled {
		return nil
	}
	tlsCfg, err := ServerTLSConfig(cfg)
	if err != nil {
		return fmt.Errorf("mtls configure server: %w", err)
	}
	srv.TLSConfig = tlsCfg
	return nil
}

// ListenConfigured serves with TLS when MTLS_ENABLED and srv.TLSConfig is set.
func ListenConfigured(srv *http.Server) error {
	cfg := LoadFromEnv()
	if cfg.Enabled {
		return srv.ListenAndServeTLS(cfg.ServerCertFile, cfg.ServerKeyFile)
	}
	return srv.ListenAndServe()
}
