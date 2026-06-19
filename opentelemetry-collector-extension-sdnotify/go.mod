// Top-level module exists only so that `go tool -modfile=internal/tools/go.mod`
// can be invoked from the repo root. The actual collector binary is generated
// by ocb under _build/ with its own go.mod.
//
// The custom extension lives at extension/sdnotify and has its own go.mod.
module github.com/iypetrov/opentelemetry-collector-extension-sdnotify

go 1.24
