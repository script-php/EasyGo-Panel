package actions

import (
	"fmt"
)

// WebServerAction handles web server operations
type WebServerAction struct {
	BaseAction
}

// NewWebServerAction creates a new web server action instance
func NewWebServerAction() *WebServerAction {
	return &WebServerAction{}
}

// InstallApache installs Apache web server
func (w *WebServerAction) InstallApache() *Result {
	// Detect package manager
	if w.FileExists("/usr/bin/apt") {
		return w.installApacheDebian()
	} else if w.FileExists("/usr/bin/yum") || w.FileExists("/usr/bin/dnf") {
		return w.installApacheRHEL()
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported Linux distribution",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}

// InstallNginx installs Nginx web server
func (w *WebServerAction) InstallNginx() *Result {
	if w.FileExists("/usr/bin/apt") {
		return w.installNginxDebian()
	} else if w.FileExists("/usr/bin/yum") || w.FileExists("/usr/bin/dnf") {
		return w.installNginxRHEL()
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported Linux distribution",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}

// ConfigureApacheVhost creates an Apache virtual host
func (w *WebServerAction) ConfigureApacheVhost(domain, docroot string) *Result {
	vhostConfig := fmt.Sprintf(`<VirtualHost *:80>
    ServerName %s
    ServerAlias www.%s
    DocumentRoot %s
    
    <Directory %s>
        Options -Indexes +FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>
    
    ErrorLog ${APACHE_LOG_DIR}/%s_error.log
    CustomLog ${APACHE_LOG_DIR}/%s_access.log combined
</VirtualHost>`, domain, domain, docroot, docroot, domain, domain)

	configPath := fmt.Sprintf("/etc/apache2/sites-available/%s.conf", domain)
	result := w.WriteFile(configPath, vhostConfig)
	if !result.Success {
		return result
	}
	
	// Enable site
	enableResult := w.RunCommand("a2ensite", domain)
	if !enableResult.Success {
		return enableResult
	}
	
	// Reload Apache
	return w.ReloadService("apache2")
}

// ConfigureNginxVhost creates an Nginx virtual host
func (w *WebServerAction) ConfigureNginxVhost(domain, docroot string) *Result {
	vhostConfig := fmt.Sprintf(`server {
    listen 80;
    server_name %s www.%s;
    root %s;
    index index.php index.html index.htm;
    
    location / {
        try_files $uri $uri/ =404;
    }
    
    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php-fpm.sock;
    }
    
    location ~ /\.ht {
        deny all;
    }
    
    access_log /var/log/nginx/%s_access.log;
    error_log /var/log/nginx/%s_error.log;
}`, domain, domain, docroot, domain, domain)

	configPath := fmt.Sprintf("/etc/nginx/sites-available/%s", domain)
	result := w.WriteFile(configPath, vhostConfig)
	if !result.Success {
		return result
	}
	
	// Enable site
	symlinkResult := w.RunCommand("ln", "-sf", configPath, fmt.Sprintf("/etc/nginx/sites-enabled/%s", domain))
	if !symlinkResult.Success {
		return symlinkResult
	}
	
	// Test configuration
	testResult := w.RunCommand("nginx", "-t")
	if !testResult.Success {
		return testResult
	}
	
	// Reload Nginx
	return w.ReloadService("nginx")
}

// Private helper methods

func (w *WebServerAction) installApacheDebian() *Result {
	updateResult := w.RunCommand("apt", "update")
	if !updateResult.Success {
		return updateResult
	}
	
	installResult := w.RunCommand("apt", "install", "-y", "apache2")
	if !installResult.Success {
		return installResult
	}
	
	// Enable and start Apache
	w.EnableService("apache2")
	return w.StartService("apache2")
}

func (w *WebServerAction) installApacheRHEL() *Result {
	var installCmd string
	if w.FileExists("/usr/bin/dnf") {
		installCmd = "dnf"
	} else {
		installCmd = "yum"
	}
	
	installResult := w.RunCommand(installCmd, "install", "-y", "httpd")
	if !installResult.Success {
		return installResult
	}
	
	// Enable and start Apache (httpd)
	w.EnableService("httpd")
	return w.StartService("httpd")
}

func (w *WebServerAction) installNginxDebian() *Result {
	updateResult := w.RunCommand("apt", "update")
	if !updateResult.Success {
		return updateResult
	}
	
	installResult := w.RunCommand("apt", "install", "-y", "nginx")
	if !installResult.Success {
		return installResult
	}
	
	// Enable and start Nginx
	w.EnableService("nginx")
	return w.StartService("nginx")
}

func (w *WebServerAction) installNginxRHEL() *Result {
	var installCmd string
	if w.FileExists("/usr/bin/dnf") {
		installCmd = "dnf"
	} else {
		installCmd = "yum"
	}
	
	installResult := w.RunCommand(installCmd, "install", "-y", "nginx")
	if !installResult.Success {
		return installResult
	}
	
	// Enable and start Nginx
	w.EnableService("nginx")
	return w.StartService("nginx")
}