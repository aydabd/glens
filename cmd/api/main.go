package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"glens/pkg/logging"
	"glens/tools/api/internal/handler"
	"glens/tools/api/internal/middleware"
)

// version is set at build time via -ldflags="-X main.version=<tag>".
var version = "dev"

func main() {
	logging.Setup(logging.Config{
		Level:  logging.LevelInfo,
		Format: logging.FormatJSON,
	})

	mux := http.NewServeMux()
	registerRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	wrapped := middleware.Recovery(middleware.Logging(middleware.CORS(mux)))

	log.Info().Str("port", port).Str("version", version).Msg("starting API server")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), wrapped); err != nil {
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
