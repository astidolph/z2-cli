package api

import (
	"fmt"
	"io/fs"
	"net/http"
)

type Server struct {
	addr       string
	frontendFS fs.FS
}

func NewServer(addr string, frontendFS fs.FS) *Server {
	return &Server{addr: addr, frontendFS: frontendFS}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/auth/status", handleAuthStatus)
	mux.HandleFunc("GET /api/config", handleGetConfig)
	mux.HandleFunc("PUT /api/config", handlePutConfig)
	mux.HandleFunc("GET /api/runs", handleGetRuns)
	mux.HandleFunc("GET /api/chart-data", handleGetChartData)
	mux.HandleFunc("POST /api/refresh", handleRefresh)

	// Serve embedded frontend if available
	if s.frontendFS != nil {
		mux.Handle("/", spaHandler(s.frontendFS))
	}

	handler := cors(mux)

	fmt.Printf("z2-cli API server running at http://localhost%s\n", s.addr)
	return http.ListenAndServe(s.addr, handler)
}

// spaHandler serves static files from the filesystem, falling back to
// index.html for any path that doesn't match a file (SPA routing).
func spaHandler(frontendFS fs.FS) http.Handler {
	fileServer := http.FileServerFS(frontendFS)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to open the requested file
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = path[1:] // strip leading slash
		}

		if _, err := fs.Stat(frontendFS, path); err != nil {
			// File not found — serve index.html for SPA routing
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})
}
