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
	if err := InitSessionKey(); err != nil {
		return err
	}

	mux := http.NewServeMux()

	// Public routes (no auth required)
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/auth/status", handleAuthStatus)
	mux.HandleFunc("GET /api/auth/login", handleAuthLogin)
	mux.HandleFunc("GET /api/auth/callback", handleAuthCallback)

	// Protected routes (session required)
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/config", handleGetConfig)
	protected.HandleFunc("PUT /api/config", handlePutConfig)
	protected.HandleFunc("GET /api/runs", handleGetRuns)
	protected.HandleFunc("GET /api/chart-data", handleGetChartData)
	protected.HandleFunc("POST /api/refresh", handleRefresh)
	protected.HandleFunc("GET /api/leaderboard", handleGetLeaderboard)
	protected.HandleFunc("POST /api/leaderboard/refresh", handleRefreshLeaderboard)
	mux.Handle("/api/", requireAuth(protected))

	// Serve embedded frontend if available
	if s.frontendFS != nil {
		mux.Handle("/", spaHandler(s.frontendFS))
	}

	handler := securityHeaders(cors(mux))

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
