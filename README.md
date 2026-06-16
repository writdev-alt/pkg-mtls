# github.com/writdev-alt/pkg-mtls

Shared mutual TLS helpers for IlonaPay internal HTTPS. Used by **api-gateway** (outbound client cert) and **production services** (inbound server with client verification).

When `MTLS_ENABLED` is off, helpers fall back to plain HTTP so local development stays unchanged.

## Architecture

```
api-gateway                    production service
(client cert: gateway.crt) --> (server.crt + verify client)
         |                              |
         +-------- shared CA (ca.crt) ---+
```

- **Gateway** — `NewTransport` / `ClientTLSConfig` present `gateway.crt` to upstreams.
- **Backends** — `ServerTLSConfig` / `RunGin` listen on HTTPS and require a valid gateway client cert.

## API

| Function | Role |
|----------|------|
| `LoadFromEnv()` | Read `MTLS_*` environment variables into `Config` |
| `ServerTLSConfig(cfg)` | Backend TLS: server cert + `RequireAndVerifyClientCert` |
| `ClientTLSConfig(cfg)` | Gateway TLS: client cert + CA trust pool |
| `NewTransport(cfg)` | `http.Transport` for gateway upstream proxy |
| `CloneTLSConfig(cfg)` | Shallow copy of `tls.Config` for per-connection use |
| `ListenAndServe(handler, addr)` | Start plain HTTP or mTLS HTTPS from env |
| `RunGin(engine, addr)` | Gin drop-in for `engine.Run` with optional mTLS |
| `ConfigureServer(srv)` | Apply server TLS settings to an existing `http.Server` |
| `ListenConfigured(srv)` | Serve using `srv.TLSConfig` when mTLS is enabled |

## Quick start

### Backend (Gin)

```go
import "github.com/writdev-alt/pkg-mtls"

if err := mtls.RunGin(router, ":"+port); err != nil {
    log.Fatal(err)
}
```

### Backend (stdlib `http.Server`)

```go
srv := &http.Server{Addr: ":8080", Handler: mux}
_ = mtls.ConfigureServer(srv)
_ = mtls.ListenConfigured(srv)
```

### Gateway (reverse proxy transport)

```go
cfg := mtls.LoadFromEnv()
transport, err := mtls.NewTransport(cfg)
if err != nil {
    log.Fatal(err)
}
client := &http.Client{Transport: transport}
```

## Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `MTLS_ENABLED` | `false` | Enable mTLS (`true`, `1`, `yes`, `on`) |
| `MTLS_CA_CERT` | `/etc/mtls/ca.crt` | Shared CA bundle |
| `MTLS_SERVER_CERT` | `/etc/mtls/server.crt` | Backend server certificate |
| `MTLS_SERVER_KEY` | `/etc/mtls/server.key` | Backend server private key |
| `MTLS_CLIENT_CERT` | `/etc/mtls/gateway.crt` | Gateway client certificate |
| `MTLS_CLIENT_KEY` | `/etc/mtls/gateway.key` | Gateway client private key |

With mTLS enabled in Compose, certs are mounted read-only at `/etc/mtls` on gateway and backend containers. Set `UPSTREAM_SCHEME=https` on the gateway so upstream URLs use HTTPS.

## Consumers in this repo

| Service | Usage |
|---------|--------|
| `api-gateway` | `NewTransport` for upstream proxy |
| `wallet-service`, `bank-service`, `transaction-service`, `worker-webhook` | `RunGin` |
| `qris-service`, `withdraw-service` | `ConfigureServer` + `ListenConfigured` |

## Tests

```bash
go test ./...
```

Unit tests cover config loading, TLS config construction, transport creation, and `ConfigureServer`. Cert fixtures are generated in-test (no checked-in PEM files).

