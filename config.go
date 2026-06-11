package mtls

import (
	"os"
	"strings"
)

// Config holds mTLS file paths and enable flag (shared by gateway and production services).
type Config struct {
	Enabled        bool
	CACertFile     string
	ServerCertFile string
	ServerKeyFile  string
	ClientCertFile string
	ClientKeyFile  string
}

// LoadFromEnv reads MTLS_* environment variables.
func LoadFromEnv() Config {
	return Config{
		Enabled:        envBool("MTLS_ENABLED", false),
		CACertFile:     env("MTLS_CA_CERT", "/etc/mtls/ca.crt"),
		ServerCertFile: env("MTLS_SERVER_CERT", "/etc/mtls/server.crt"),
		ServerKeyFile:  env("MTLS_SERVER_KEY", "/etc/mtls/server.key"),
		ClientCertFile: env("MTLS_CLIENT_CERT", "/etc/mtls/gateway.crt"),
		ClientKeyFile:  env("MTLS_CLIENT_KEY", "/etc/mtls/gateway.key"),
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
