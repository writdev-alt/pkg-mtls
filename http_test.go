package mtls

import (
	"net/http"
	"strings"
	"testing"
)

func TestConfigureServer_disabled(t *testing.T) {
	t.Setenv("MTLS_ENABLED", "false")

	srv := &http.Server{Addr: ":0"}
	if err := ConfigureServer(srv); err != nil {
		t.Fatalf("ConfigureServer() err = %v", err)
	}
	if srv.TLSConfig != nil {
		t.Fatal("expected TLSConfig unset when disabled")
	}
}

func TestConfigureServer_enabled(t *testing.T) {
	cfg := testConfig(t)
	t.Setenv("MTLS_ENABLED", "true")
	t.Setenv("MTLS_CA_CERT", cfg.CACertFile)
	t.Setenv("MTLS_SERVER_CERT", cfg.ServerCertFile)
	t.Setenv("MTLS_SERVER_KEY", cfg.ServerKeyFile)
	t.Setenv("MTLS_CLIENT_CERT", cfg.ClientCertFile)
	t.Setenv("MTLS_CLIENT_KEY", cfg.ClientKeyFile)

	srv := &http.Server{Addr: ":0"}
	if err := ConfigureServer(srv); err != nil {
		t.Fatalf("ConfigureServer() err = %v", err)
	}
	if srv.TLSConfig == nil {
		t.Fatal("expected TLSConfig when enabled")
	}
}

func TestConfigureServer_enabled_missingCert(t *testing.T) {
	t.Setenv("MTLS_ENABLED", "true")
	t.Setenv("MTLS_CA_CERT", t.TempDir()+"/missing-ca.crt")
	t.Setenv("MTLS_SERVER_CERT", "/etc/mtls/server.crt")
	t.Setenv("MTLS_SERVER_KEY", "/etc/mtls/server.key")

	srv := &http.Server{Addr: ":0"}
	err := ConfigureServer(srv)
	if err == nil || !strings.Contains(err.Error(), "mtls configure server") {
		t.Fatalf("ConfigureServer() err = %v", err)
	}
}
