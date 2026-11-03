# Roadmap

Engi is a lightweight Go web framework built around **services** — a unit of code owns a
route prefix, a set of routes, and the parameters those routes accept. The framework
handles routing, parameter validation, response wrapping, and (planned) API documentation
so that application code stays declarative.

This document describes the work planned for the project. Items are grouped by release and
carry an effort estimate so the scope is explicit. Every item is tracked as a GitHub
issue — the links are given inline.

**Status legend:** `🔜 planned` · `🚧 in progress` · `✅ done`

---

## Current state

| Area | Status |
|---|---|
| Trie-based router (`{param}`, `*wildcard`) | ✅ done |
| Service registration, prefixes, graceful shutdown | ✅ done |
| Parameter extraction (path / query / header / cookie) with type parsing | ✅ done |
| Parameter validators (`NotEmpty`, `Greater`, `Less`, `AND`, `OR`) | ✅ done |
| Middleware pipeline with priority ordering | ✅ done |
| Response wrappers (`ResponseAsIs`, `ResponseAsObject`), JSON / XML marshaling | ✅ done |
| CORS middleware | ✅ done |
| Auth middleware (Basic, Bearer, API key) | 🚧 partially broken — see v0.2 |
| OpenTelemetry tracing (via `otelhttp`) | ✅ done |
| OpenAPI / Swagger generation | 🔜 planned — see v0.3 |
| Test coverage | 🚧 thin — see v0.3 |
| Benchmarks | 🔜 planned — see v0.3 |

---

## v0.2 — Correctness and stability

**Goal:** fix the defects listed below and make error handling consistent across the
framework. Estimated effort: **~8 hours**.

### Bugs to fix

- [x] **Routes with a path parameter never matched.** ([#26](https://github.com/kliuchnikovv/engi/issues/26))
      The trie router recognized `:name` segments, but the public API — the README,
      the example application and every `path.Integer("id", …)` declaration — uses the
      `{name}` syntax. The trie was therefore matching the literal segment `{id}` and
      answering `404` for every route that had a path parameter. The trie now accepts
      both spellings.
- [x] **The trailing wildcard syntax was audited** and is covered by tests.
      `GET /assets/*filepath` matches through the public API and hands the handler the
      swallowed slashes (`css/main.css`). The example application and the documentation
      use only the `{name}` syntax, and no route outside the tests declares a wildcard.
- [ ] **Bearer auth never validates the presented token.** ([#10](https://github.com/kliuchnikovv/engi/issues/10))
      `definition/middlewares/auth/auth.go:98` compares the constant `authHeader`
      (`"Authorization"`) instead of the value read from the request, and the branch is
      inverted — a *valid* token is rejected while an *invalid* one is accepted. Bearer
      auth is effectively a no-op for any request carrying an `Authorization` header.
      Add regression tests covering valid / invalid / missing / malformed tokens.
- [ ] **Handler errors are swallowed.** ([#11](https://github.com/kliuchnikovv/engi/issues/11))
      `internal/routes/routes.go:79` prints the error with `fmt.Print` and returns `nil`,
      so a failing handler produces a `200 OK` with an empty body. Errors must reach the
      response writer and the logger.
- [ ] **Every middleware failure is reported as `400 Bad Request`.** ([#12](https://github.com/kliuchnikovv/engi/issues/12))
      `internal/routes/route.go:78` maps all middleware errors through
      `response.BadRequest`, so an auth failure answers `400` instead of `401`, and a
      rate-limit failure answers `400` instead of `429`. Introduce a typed error that
      middlewares can use to select the status code.
- [ ] **Parameter conversion panics instead of returning an error.** ([#13](https://github.com/kliuchnikovv/engi/issues/13))
      `internal/request/request.go` panics in `Bool`, `Integer`, `Float`, `String` and
      `Time` when a parsed value has an unexpected type. Replace panics with errors
      surfaced through the response.
- [ ] **`Engine.Shutdown` ignores its context.** ([#14](https://github.com/kliuchnikovv/engi/issues/14))
      `server.go:84` discards the `context.Context` argument and only signals the
      internal channel; the deadline is never honoured.
- [ ] **API key / basic auth compare secrets with `!=`.** ([#15](https://github.com/kliuchnikovv/engi/issues/15))
      `definition/middlewares/auth/auth.go` — switch to `crypto/subtle` constant-time
      comparison so token checks do not leak length or prefix information through
      response timing.
- [ ] **API key query parameters are joined with an empty separator.**
      `definition/middlewares/auth/auth.go:125` joins repeated query values with `""`;
      use `","` for consistency with `Request.GetParameter`.

### Cleanup

- [ ] Remove dead code left over from the regexp router: `parameterRegexp`
      (`internal/routes/routes.go:18`) and the `regexpRoute` type are only referenced from
      commented-out blocks. ([#16](https://github.com/kliuchnikovv/engi/issues/16))
- [ ] Remove the discarded `slog.New(...)` call in `NewService`
      (`service.go:58`), whose result is thrown away.
- [ ] Remove the duplicated `Start` doc comment (`server.go:58` and `server.go:62`).
- [ ] Call `signal.Stop` on shutdown so the process stops intercepting `SIGTERM`.
- [ ] The remaining small defects. ([#17](https://github.com/kliuchnikovv/engi/issues/17))

---

## v0.3 — Documentation, testing, performance

**Goal:** make the framework self-describing and prove it behaves correctly.
Estimated effort: **~14 hours**.

### OpenAPI generation

The `internal/docs` package exists but is empty, and every `Middleware.Docs(route)`
implementation panics with `not implemented`. The pieces for a real specification are
already in place — routes carry a description, parameters carry a type name and a
validation regexp, and responses are marshaled by a pluggable `Marshaler`.

- [ ] Define the OpenAPI 3.0 document model in `internal/docs/types.go`
      (paths, operations, parameters, request bodies, responses, security schemes).
- [ ] Implement `Docs(route *routes.Route)` for every middleware type:
      descriptions, parameter schemas, auth requirements.
- [ ] Collect the specification from the registered services.
- [ ] Serve the generated document over HTTP.
- [ ] Support per-service or per-route overrides for tags, summaries and status codes.

Tracked as [#18](https://github.com/kliuchnikovv/engi/issues/18).

- [ ] Serve a Swagger UI / Redoc page pointed at that document.

Tracked as [#19](https://github.com/kliuchnikovv/engi/issues/19).

### Tests

- [ ] Table-driven tests for the trie router: exact, parameterized, wildcard, method
      mismatch, and route precedence between static and parameterized segments.
- [ ] Tests for parameter extraction and every validator in `definition/validate`.
- [ ] Tests for the response wrappers and both marshaling formats.
- [ ] Tests for the middleware pipeline: priority ordering, short-circuiting on error,
      and per-service vs per-route middleware composition.
- [ ] Integration tests that drive a full engine over `httptest`.
- [ ] Raise coverage to a level worth publishing as a badge.

Tracked as [#20](https://github.com/kliuchnikovv/engi/issues/20).

### Benchmarks

- [ ] Benchmarks for route matching across static, parameterized and wildcard paths.
- [ ] Benchmarks for parameter parsing and validation.
- [ ] Benchmarks for response marshaling in both formats.

Tracked as [#21](https://github.com/kliuchnikovv/engi/issues/21).

### Developer experience

- [ ] Replace the `fmt.Print` debugging leftovers with structured logging.
      ([#11](https://github.com/kliuchnikovv/engi/issues/11))
- [ ] Move the hard-coded middleware priorities into an exported map so applications can
      order their own middleware deterministically.
      ([#24](https://github.com/kliuchnikovv/engi/issues/24))

---

## v0.4 — Extended capabilities

**Goal:** cover the transports and middleware that a production service usually needs.
Estimated effort: **~16 hours**.

- [ ] **gRPC transport.** Expose services over gRPC alongside HTTP, mapping parameters and
      responses onto the protobuf types (`internal/response/response.go`).
      ([#22](https://github.com/kliuchnikovv/engi/issues/22))
- [ ] **WebSocket and server-sent events.** Support upgrading a route to a streaming
      handler, and add an SSE response helper.
      ([#23](https://github.com/kliuchnikovv/engi/issues/23))
- [ ] **Rate limiting middleware.** Token bucket per client, per route or per service.
      ([#24](https://github.com/kliuchnikovv/engi/issues/24))
- [ ] **Request/response compression.** Gzip and Brotli content encodings.
      ([#24](https://github.com/kliuchnikovv/engi/issues/24))
- [ ] **Static file serving.** Serve embedded assets or a directory from a route.
- [ ] **Panic recovery middleware.** Convert handler panics into `500` responses instead
      of dropping the connection. ([#24](https://github.com/kliuchnikovv/engi/issues/24))
- [ ] **Context helpers.** Propagate values through `context.Context` rather than through
      the request wrapper, so handlers can carry request-scoped state.

---

## v0.5 — Ecosystem and release engineering

**Goal:** make the framework pleasant to adopt and to maintain.

- [ ] Publish versioned releases with a changelog and semantic versioning.
- [ ] Add a `CONTRIBUTING.md` and a code of conduct.
- [ ] Expand the example application in `kliuchnikovv/engi-example`.
- [ ] Add a middleware gallery in the documentation — each middleware with a runnable
      snippet and its configuration options.
- [ ] CI: run the test suite against the supported Go versions, publish coverage, and
      generate release artifacts.
- [ ] Package documentation for every exported identifier, with usage examples that
      compile as documentation examples.
- [ ] Enable request body size limits (`engi.go:14`).
- [ ] Publish request/response access logs with duration and status (`engi.go:17`).

Tracked as [#25](https://github.com/kliuchnikovv/engi/issues/25).

---

## Principles

These guide the work above and any change that is not yet listed.

1. **The framework removes boilerplate; it does not hide the standard library.** Handlers
   receive a real `context.Context` and a real `*http.Request`.
2. **Declarative over imperative.** A route should describe what it needs — parameters,
   validation, middleware, documentation — and the engine should enforce it.
3. **Explicit errors.** Failures are typed and surfaced as HTTP status codes, never
   swallowed and never panicked.
4. **No mandatory dependencies beyond the standard library and OpenTelemetry.**

---

## Contributing

Issues and pull requests are welcome. Before opening a pull request, please run:

```sh
go build ./...
go test ./...
go vet ./...
```

Work that is listed above as `🔜 planned` is tracked as a GitHub issue — please comment on
an existing issue before starting, so that effort is not duplicated.
