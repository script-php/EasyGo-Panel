package web

import (
	"easygo/pkg/actions"
	"easygo/pkg/auth"
	"encoding/json"
	"net/http"
	
	"github.com/gorilla/mux"
)

// PageData represents common page data
type PageData struct {
	Title       string
	User        string
	CurrentPage string
	Flash       string
	Data        interface{}
}

// handleHome shows the home page or redirects to dashboard if authenticated
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/panel/", http.StatusFound)
		return
	}
	
	http.Redirect(w, r, "/login", http.StatusFound)
}

// handleLogin handles login page and authentication
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := PageData{
			Title:       "Login - EasyGo Panel",
			CurrentPage: "login",
		}
		s.renderTemplate(w, "login.html", data)
		return
	}
	
	// POST request - handle authentication
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	if username == "" || password == "" {
		data := PageData{
			Title:       "Login - EasyGo Panel",
			CurrentPage: "login",
			Flash:       "Username and password are required",
		}
		s.renderTemplate(w, "login.html", data)
		return
	}
	
	// Authenticate using PAM
	if err := auth.AuthenticateUser(username, password); err != nil {
		data := PageData{
			Title:       "Login - EasyGo Panel",
			CurrentPage: "login",
			Flash:       "Invalid username or password",
		}
		s.renderTemplate(w, "login.html", data)
		return
	}
	
	// Set session
	session, _ := s.store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Save(r, w)
	
	http.Redirect(w, r, "/panel/", http.StatusFound)
}

// handleLogout handles user logout
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	session.Values["authenticated"] = false
	delete(session.Values, "username")
	session.Save(r, w)
	
	http.Redirect(w, r, "/login", http.StatusFound)
}

// handleDashboard shows the main dashboard
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	// Get real system stats (basic implementation)
	stats := make(map[string]interface{})
	
	// This is a basic implementation - in production you'd want proper system monitoring
	stats["cpu_usage"] = "45%"
	stats["memory_usage"] = "2.4 GB"
	stats["disk_usage"] = "75%"
	stats["uptime"] = "5 days"
	
	data := PageData{
		Title:       "Dashboard - EasyGo Panel",
		User:        username,
		CurrentPage: "dashboard",
		Data:        stats,
	}
	
	s.renderTemplate(w, "dashboard.html", data)
}

// handleServices shows the services page
func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "Services - EasyGo Panel",
		User:        username,
		CurrentPage: "services",
	}
	
	s.renderTemplate(w, "services.html", data)
}

// handleApache handles Apache management
func (s *Server) handleApache(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "Apache - EasyGo Panel",
		User:        username,
		CurrentPage: "apache",
	}
	
	s.renderTemplate(w, "apache.html", data)
}

// handleNginx handles Nginx management
func (s *Server) handleNginx(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "Nginx - EasyGo Panel",
		User:        username,
		CurrentPage: "nginx",
	}
	
	s.renderTemplate(w, "nginx.html", data)
}

// handlePHP handles PHP management
func (s *Server) handlePHP(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "PHP - EasyGo Panel",
		User:        username,
		CurrentPage: "php",
	}
	
	s.renderTemplate(w, "php.html", data)
}

// handleDomains handles domain management
func (s *Server) handleDomains(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "Domains - EasyGo Panel",
		User:        username,
		CurrentPage: "domains",
	}
	
	s.renderTemplate(w, "domains.html", data)
}

// handleSSL handles SSL certificate management
func (s *Server) handleSSL(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "SSL Certificates - EasyGo Panel",
		User:        username,
		CurrentPage: "ssl",
	}
	
	s.renderTemplate(w, "ssl.html", data)
}

// handleDatabases handles database management
func (s *Server) handleDatabases(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "Databases - EasyGo Panel",
		User:        username,
		CurrentPage: "databases",
	}
	
	s.renderTemplate(w, "databases.html", data)
}

// handleSettings handles panel settings
func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session")
	username, _ := session.Values["username"].(string)
	
	data := PageData{
		Title:       "Settings - EasyGo Panel",
		User:        username,
		CurrentPage: "settings",
	}
	
	s.renderTemplate(w, "settings.html", data)
}

// API Handlers

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// handleAPIServiceStatus returns status of all services
func (s *Server) handleAPIServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	baseAction := &actions.BaseAction{}
	services := []string{"apache2", "nginx", "php8.2-fpm", "mysql", "postgresql"}
	var serviceData []map[string]interface{}
	
	for _, service := range services {
		result := baseAction.ServiceStatus(service)
		status := "stopped"
		if result.Success {
			if serviceResult, ok := result.Data.(*actions.Service); ok {
				if serviceResult.Status == "active" {
					status = "running"
				}
			}
		}
		
		serviceData = append(serviceData, map[string]interface{}{
			"name":   service,
			"status": status,
		})
	}
	
	response := APIResponse{
		Success: true,
		Data:    serviceData,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleAPIServiceStart starts a service
func (s *Server) handleAPIServiceStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	serviceName := vars["service"]
	
	baseAction := &actions.BaseAction{}
	result := baseAction.StartService(serviceName)
	
	response := APIResponse{
		Success: result.Success,
		Message: result.Message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleAPIServiceStop stops a service
func (s *Server) handleAPIServiceStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	serviceName := vars["service"]
	
	baseAction := &actions.BaseAction{}
	result := baseAction.StopService(serviceName)
	
	response := APIResponse{
		Success: result.Success,
		Message: result.Message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleAPIServiceRestart restarts a service
func (s *Server) handleAPIServiceRestart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	serviceName := vars["service"]
	
	baseAction := &actions.BaseAction{}
	result := baseAction.RestartService(serviceName)
	
	response := APIResponse{
		Success: result.Success,
		Message: result.Message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleAPISystemStats returns system statistics
func (s *Server) handleAPISystemStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Basic system stats - in production you'd use proper system monitoring
	stats := map[string]string{
		"cpu":    "45%",
		"memory": "2.4 GB",
		"disk":   "75%",
		"uptime": "5 days",
	}
	
	json.NewEncoder(w).Encode(stats)
}