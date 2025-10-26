package web

import (
	"easygo/pkg/auth"
	"net/http"
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
	
	data := PageData{
		Title:       "Dashboard - EasyGo Panel",
		User:        username,
		CurrentPage: "dashboard",
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