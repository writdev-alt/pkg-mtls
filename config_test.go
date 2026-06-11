package mtls

import "testing"

func TestLoadFromEnv_defaults(t *testing.T) {
	t.Setenv("MTLS_ENABLED", "")
	t.Setenv("MTLS_CA_CERT", "")
	t.Setenv("MTLS_SERVER_CERT", "")
	t.Setenv("MTLS_SERVER_KEY", "")
	t.Setenv("MTLS_CLIENT_CERT", "")
	t.Setenv("MTLS_CLIENT_KEY", "")

	cfg := LoadFromEnv()
	if cfg.Enabled {
		t.Fatal("expected MTLS disabled by default")
	}
	if cfg.CACertFile != "/etc/mtls/ca.crt" {
		t.Fatalf("CACertFile = %q", cfg.CACertFile)
	}
	if cfg.ServerCertFile != "/etc/mtls/server.crt" {
		t.Fatalf("ServerCertFile = %q", cfg.ServerCertFile)
	}
	if cfg.ClientCertFile != "/etc/mtls/gateway.crt" {
		t.Fatalf("ClientCertFile = %q", cfg.ClientCertFile)
	}
}

func TestLoadFromEnv_overrides(t *testing.T) {
	t.Setenv("MTLS_ENABLED", "true")
	t.Setenv("MTLS_CA_CERT", "/custom/ca.pem")
	t.Setenv("MTLS_SERVER_CERT", "/custom/server.crt")
	t.Setenv("MTLS_SERVER_KEY", "/custom/server.key")
	t.Setenv("MTLS_CLIENT_CERT", "/custom/client.crt")
	t.Setenv("MTLS_CLIENT_KEY", "/custom/client.key")

	cfg := LoadFromEnv()
	if !cfg.Enabled {
		t.Fatal("expected MTLS enabled")
	}
	if cfg.CACertFile != "/custom/ca.pem" {
		t.Fatalf("CACertFile = %q", cfg.CACertFile)
	}
	if cfg.ClientKeyFile != "/custom/client.key" {
		t.Fatalf("ClientKeyFile = %q", cfg.ClientKeyFile)
	}
}

func TestEnvBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		fallback bool
		want     bool
	}{
		{"empty uses fallback true", "", true, true},
		{"empty uses fallback false", "", false, false},
		{"true", "true", false, true},
		{"1", "1", false, true},
		{"yes", "yes", false, true},
		{"on", "ON", false, true},
		{"false", "false", true, false},
		{"0", "0", true, false},
		{"invalid uses fallback", "maybe", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MTLS_TEST_BOOL", tt.value)
			if got := envBool("MTLS_TEST_BOOL", tt.fallback); got != tt.want {
				t.Fatalf("envBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
