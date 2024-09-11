package httphandler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"path"

	"github.com/vizee/export"
)

type Handler struct{}

func (*Handler) List(w http.ResponseWriter, r *http.Request) {
	varNames := export.AllVars()
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	err := enc.Encode(varNames)
	if err != nil {
		slog.Error("json encode", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (*Handler) Get(w http.ResponseWriter, r *http.Request) {
	v := export.Lookup(r.PathValue("name"))
	if v == nil {
		http.Error(w, "No such key", http.StatusNotFound)
		return
	}
	_, _ = w.Write(v.Get())
}

func (*Handler) Set(w http.ResponseWriter, r *http.Request) {
	v := export.Lookup(r.PathValue("name"))
	if v == nil {
		http.Error(w, "No such key", http.StatusNotFound)
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("read body", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = v.Set(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}

func AddMux(mux *http.ServeMux, prefix string) {
	h := &Handler{}
	mux.Handle("GET "+path.Join(prefix, "/vars"), http.HandlerFunc(h.List))
	mux.Handle("GET "+path.Join(prefix, "/var/{name}"), http.HandlerFunc(h.Get))
	mux.Handle("PUT "+path.Join(prefix, "/var/{name}"), http.HandlerFunc(h.Set))
}
