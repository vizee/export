package main

import (
	"log/slog"
	"net/http"

	"github.com/vizee/export"
	_ "github.com/vizee/export/debug"
	"github.com/vizee/export/httphandler"
)

var (
	testInt int = 1
)

func main() {
	export.Int("test", &testInt)

	http.HandleFunc("GET /print", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("print", "test", testInt)
	})

	httphandler.AddMux(http.DefaultServeMux, "/_")

	slog.Error("listen http", "err", http.ListenAndServe(":8080", nil))
}
