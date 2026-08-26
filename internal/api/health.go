package api

import (
	"io"
	"net/http"

	"frame-static/internal/assemble"
	"frame-static/internal/buckling"
	"frame-static/internal/model"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, APIError{Code: "method", Message: "GET required"})
		return
	}
	writeJSON(w, map[string]any{"ok": true, "service": "frame-static"})
}

func BucklingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, APIError{Code: "method", Message: "POST required"})
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBody))
	if err != nil {
		writeError(w, APIError{Code: "read", Message: err.Error()})
		return
	}
	m, err := model.ParseModelBytes(body)
	if err != nil {
		writeError(w, Classify(err))
		return
	}
	res, err := assemble.Solve(m)
	if err != nil {
		writeError(w, Classify(err))
		return
	}
	checks, err := buckling.Frame(m, res, buckling.PinnedK())
	if err != nil {
		writeError(w, APIError{Code: "buckling", Message: err.Error()})
		return
	}
	worst, _ := buckling.Worst(checks)
	writeJSON(w, map[string]any{
		"ok":     true,
		"checks": checks,
		"worst":  worst,
	})
}
