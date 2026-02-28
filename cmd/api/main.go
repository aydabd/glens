package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"glens/pkg/logging"
	"glens/tools/api/internal/handler"
	"glens/tools/api/internal/middleware"
)

// version is set at build time via -ldflags="-X main.version=<tag>".
var version = "dev"

func main() {
	level := logging.LevelInfo
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		level = logging.Level(envLevel)
	}

	logging.Setup(logging.Config{
		Level:  level,
		Format: logging.FormatJSON,
	})

	mux := http.NewServeMux()
	registerRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	wrapped := middleware.Recovery(middleware.Logging(middleware.CORS(mux)))

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           wrapped,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Info().Str("port", port).Str("version", version).Msg("starting API server")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", handler.Health(version))
	mux.HandleFunc("POST /api/v1/analyze", handler.Analyze)
	mux.HandleFunc("POST /api/v1/analyze/preview", handler.AnalyzePreview)
	mux.HandleFunc("GET /api/v1/models", handler.Models)
	mux.HandleFunc("POST /api/v1/mcp", handler.MCP)
}
