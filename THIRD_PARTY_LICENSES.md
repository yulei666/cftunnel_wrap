# Third-Party Licenses

`cftunnel` statically links the Go modules below. This list was generated from
`go version -m cftunnel` for a build against cloudflared `2026.9.3`.
License types were identified from each module's LICENSE file. The full license
texts ship with each module's source; run `go mod download` and look in
`$(go env GOMODCACHE)`.

The Go standard library is distributed under the BSD-3-Clause license
(https://go.dev/LICENSE).

| Module | Version | License |
|---|---|---|
| `github.com/BurntSushi/toml` | v1.2.0 | MIT |
| `github.com/beorn7/perks` | v1.0.1 | MIT |
| `github.com/cespare/xxhash/v2` | v2.3.0 | MIT |
| `github.com/cloudflare/backoff` | v0.0.0-20240920015135-e46b80a3a7d0 | BSD-2-Clause |
| `github.com/cloudflare/cloudflared` | v0.0.0-20260924153110-96d39adbc812 | Apache-2.0 |
| `github.com/coreos/go-oidc/v3` | v3.17.0 | Apache-2.0 |
| `github.com/coreos/go-systemd/v22` | v22.5.0 | Apache-2.0 |
| `github.com/cpuguy83/go-md2man/v2` | v2.0.0 | MIT |
| `github.com/facebookgo/grace` | v0.0.0-20180706040059-75cf19382434 | MIT |
| `github.com/fsnotify/fsnotify` | v1.4.9 | BSD-3-Clause |
| `github.com/getsentry/sentry-go` | v0.43.0 | MIT |
| `github.com/go-chi/chi/v5` | v5.3.1 | MIT |
| `github.com/go-chi/cors` | v1.2.1 | MIT |
| `github.com/go-jose/go-jose/v4` | v4.1.4 | Apache-2.0 |
| `github.com/go-logr/logr` | v1.4.4 | Apache-2.0 |
| `github.com/go-logr/stdr` | v1.2.2 | Apache-2.0 |
| `github.com/gobwas/httphead` | v0.1.0 | MIT |
| `github.com/gobwas/pool` | v0.2.1 | MIT |
| `github.com/gobwas/ws` | v1.2.1 | MIT |
| `github.com/google/gopacket` | v1.1.19 | BSD-3-Clause |
| `github.com/google/uuid` | v1.6.0 | BSD-3-Clause |
| `github.com/gorilla/websocket` | v1.5.3 | BSD-2-Clause |
| `github.com/grpc-ecosystem/grpc-gateway/v2` | v2.29.0 | BSD-3-Clause |
| `github.com/json-iterator/go` | v1.1.12 | MIT |
| `github.com/klauspost/compress` | v1.18.0 | Apache-2.0 |
| `github.com/mattn/go-colorable` | v0.1.13 | MIT |
| `github.com/mattn/go-isatty` | v0.0.20 | MIT |
| `github.com/mitchellh/go-homedir` | v1.1.0 | MIT |
| `github.com/modern-go/concurrent` | v0.0.0-20180306012644-bacd9c7ef1dd | Apache-2.0 |
| `github.com/modern-go/reflect2` | v1.0.2 | Apache-2.0 |
| `github.com/munnerz/goautoneg` | v0.0.0-20191010083416-a7dc8b61c822 | BSD-3-Clause |
| `github.com/pkg/errors` | v0.9.1 | BSD-2-Clause |
| `github.com/prometheus/client_golang` | v1.22.0 | Apache-2.0 |
| `github.com/prometheus/client_model` | v0.6.2 | Apache-2.0 |
| `github.com/prometheus/common` | v0.64.0 | Apache-2.0 |
| `github.com/prometheus/procfs` | v0.15.1 | Apache-2.0 |
| `github.com/quic-go/quic-go (=> github.com/chungthuang/quic-go)` | v0.45.1-0.20260529212404-a9fddf436fc4 | MIT |
| `github.com/rs/zerolog` | v1.20.0 | MIT |
| `github.com/russross/blackfriday/v2` | v2.1.0 | BSD-2-Clause |
| `github.com/shirou/gopsutil/v4` | v4.26.3 | BSD-3-Clause |
| `github.com/tklauser/go-sysconf` | v0.3.16 | BSD-3-Clause |
| `github.com/tklauser/numcpus` | v0.11.0 | Apache-2.0 |
| `github.com/urfave/cli/v2 (=> github.com/ipostelnik/cli/v2)` | v2.3.1-0.20210324024421-b6ea8234fe3d | MIT |
| `go.opentelemetry.io/auto/sdk` | v1.2.1 | Apache-2.0 |
| `go.opentelemetry.io/contrib/propagators` | v0.22.0 | Apache-2.0 |
| `go.opentelemetry.io/otel` | v1.45.0 | Apache-2.0 |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace` | v1.45.0 | Apache-2.0 |
| `go.opentelemetry.io/otel/metric` | v1.45.0 | Apache-2.0 |
| `go.opentelemetry.io/otel/sdk` | v1.45.0 | Apache-2.0 |
| `go.opentelemetry.io/otel/trace` | v1.45.0 | Apache-2.0 |
| `go.opentelemetry.io/proto/otlp` | v1.11.0 | Apache-2.0 |
| `golang.org/x/crypto` | v0.55.0 | BSD-3-Clause |
| `golang.org/x/net` | v0.58.0 | BSD-3-Clause |
| `golang.org/x/oauth2` | v0.36.0 | BSD-3-Clause |
| `golang.org/x/sync` | v0.22.0 | BSD-3-Clause |
| `golang.org/x/sys` | v0.47.0 | BSD-3-Clause |
| `golang.org/x/term` | v0.45.0 | BSD-3-Clause |
| `golang.org/x/text` | v0.41.0 | BSD-3-Clause |
| `google.golang.org/genproto/googleapis/api` | v0.0.0-20260720211330-0afa2a65878a | Apache-2.0 |
| `google.golang.org/genproto/googleapis/rpc` | v0.0.0-20260720211330-0afa2a65878a | Apache-2.0 |
| `google.golang.org/grpc` | v1.83.2 | Apache-2.0 |
| `google.golang.org/protobuf` | v1.36.11 | BSD-3-Clause |
| `gopkg.in/natefinch/lumberjack.v2` | v2.0.0 | MIT |
| `gopkg.in/yaml.v2` | v2.4.0 | Apache-2.0 + MIT |
| `gopkg.in/yaml.v3 (=> gopkg.in/yaml.v3)` | v3.0.1 | Apache-2.0 |
| `nhooyr.io/websocket` | v1.8.7 | MIT |
| `zombiezen.com/go/capnproto2` | v2.18.0+incompatible | MIT |
