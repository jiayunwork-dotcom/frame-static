package api

import (
	"io/fs"
	"net/http"
)

func New(webFS, exampleFS fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/health", HealthHandler)
	mux.HandleFunc("/api/solve", SolveHandler)
	mux.HandleFunc("/api/meta", MetaHandler)
	mux.HandleFunc("/api/buckling", BucklingHandler)
	mux.Handle("/", staticHandler(webFS, exampleFS))
	return withLogging(mux)
}
