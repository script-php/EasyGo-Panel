package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//go:embed all:assets
var embedFS embed.FS

// Server represents the web server
type Server struct {
	router   *mux.Router
	store    *sessions.CookieStore
	template *template.Template
}

// NewServer creates a new web server instance
func NewServer() *Server {
	server := &Server{
		router: mux.NewRouter(),
		store:  sessions.NewCookieStore([]byte("easygo-secret-key-change-this")),
	}
	
	// Parse templates
	var err error
	server.template, err = template.ParseFS(embedFS, "assets/templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
	
	server.setupRoutes()
	return server
}

// Start starts the web server
func (s *Server) Start(addr string) error {
	log.Printf("Starting EasyGo Web Panel on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// setupRoutes configures all the routes
func (s *Server) setupRoutes() {
	// Static files
	staticFS, err := fs.Sub(embedFS, "assets/static")
	if err != nil {
		log.Fatal("Failed to create static filesystem:", err)
	}
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	
	// Authentication routes
	s.router.HandleFunc("/", s.handleHome).Methods("GET")
	s.router.HandleFunc("/login", s.handleLogin).Methods("GET", "POST")
	s.router.HandleFunc("/logout", s.handleLogout).Methods("POST")
	
	// Protected routes
	protected := s.router.PathPrefix("/panel").Subrouter()
	protected.Use(s.authMiddleware)
	
	// Dashboard
	protected.HandleFunc("/", s.handleDashboard).Methods("GET")
	
	// Services
	protected.HandleFunc("/services", s.handleServices).Methods("GET")
	protected.HandleFunc("/services/apache", s.handleApache).Methods("GET", "POST")
	protected.HandleFunc("/services/nginx", s.handleNginx).Methods("GET", "POST")
	protected.HandleFunc("/services/php", s.handlePHP).Methods("GET", "POST")
	
	// Domains
	protected.HandleFunc("/domains", s.handleDomains).Methods("GET", "POST")
	
	// SSL
	protected.HandleFunc("/ssl", s.handleSSL).Methods("GET", "POST")
	
	// Databases
	protected.HandleFunc("/databases", s.handleDatabases).Methods("GET", "POST")
	
	// Settings
	protected.HandleFunc("/settings", s.handleSettings).Methods("GET", "POST")
	
	// API endpoints
	api := protected.PathPrefix("/api").Subrouter()
	api.HandleFunc("/services/status", s.handleAPIServiceStatus).Methods("GET")
	api.HandleFunc("/services/{service}/start", s.handleAPIServiceStart).Methods("POST")
	api.HandleFunc("/services/{service}/stop", s.handleAPIServiceStop).Methods("POST")
	api.HandleFunc("/services/{service}/restart", s.handleAPIServiceRestart).Methods("POST")
	api.HandleFunc("/services/{service}/uninstall", s.handleAPIServiceUninstall).Methods("POST")
	api.HandleFunc("/system/stats", s.handleAPISystemStats).Methods("GET")
}

// authMiddleware checks if user is authenticated
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "session")
		
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// renderTemplate renders a template with data
func (s *Server) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := s.template.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}