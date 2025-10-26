package actions

import (
	"fmt"
	"strings"
)

// PHPAction handles PHP operations
type PHPAction struct {
	BaseAction
}

// NewPHPAction creates a new PHP action instance
func NewPHPAction() *PHPAction {
	return &PHPAction{}
}

// PHPVersion represents a PHP version configuration
type PHPVersion struct {
	Version     string
	Installed   bool
	FPMRunning  bool
	ConfigPath  string
	FPMPath     string
}

// GetAvailableVersions returns available PHP versions
func (p *PHPAction) GetAvailableVersions() []string {
	return []string{"5.6", "7.0", "7.1", "7.2", "7.3", "7.4", "8.0", "8.1", "8.2", "8.3", "8.4"}
}

// InstallPHP installs a specific PHP version with common extensions
func (p *PHPAction) InstallPHP(version string) *Result {
	if p.FileExists("/usr/bin/apt") {
		return p.installPHPDebian(version)
	} else if p.FileExists("/usr/bin/yum") || p.FileExists("/usr/bin/dnf") {
		return p.installPHPRHEL(version)
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported Linux distribution",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}

// GetInstalledVersions returns installed PHP versions
func (p *PHPAction) GetInstalledVersions() *Result {
	var versions []*PHPVersion
	
	for _, version := range p.GetAvailableVersions() {
		phpVersion := &PHPVersion{
			Version: version,
		}
		
		// Check if PHP binary exists
		phpBinary := fmt.Sprintf("/usr/bin/php%s", version)
		
		if p.FileExists(phpBinary) {
			phpVersion.Installed = true
			phpVersion.ConfigPath = fmt.Sprintf("/etc/php/%s", version)
			phpVersion.FPMPath = fmt.Sprintf("/etc/php/%s/fpm", version)
			
			// Check if FPM is running
			serviceName := fmt.Sprintf("php%s-fpm", version)
			statusResult := p.ServiceStatus(serviceName)
			if statusResult.Success {
				if service, ok := statusResult.Data.(*Service); ok {
					phpVersion.FPMRunning = service.Status == "active"
				}
			}
		}
		
		if phpVersion.Installed {
			versions = append(versions, phpVersion)
		}
	}
	
	return &Result{
		Success: true,
		Message: fmt.Sprintf("Found %d installed PHP versions", len(versions)),
		Data:    versions,
	}
}

// ConfigurePHPFPM configures PHP-FPM for a specific version
func (p *PHPAction) ConfigurePHPFPM(version string, poolName string) *Result {
	poolConfig := fmt.Sprintf(`[%s]
user = www-data
group = www-data
listen = /var/run/php/php%s-fpm-%s.sock
listen.owner = www-data
listen.group = www-data
listen.mode = 0660

pm = dynamic
pm.max_children = 50
pm.start_servers = 5
pm.min_spare_servers = 5
pm.max_spare_servers = 35
pm.max_requests = 500

php_admin_value[sendmail_path] = /usr/sbin/sendmail -t -i -f www@localhost
php_flag[display_errors] = off
php_admin_value[error_log] = /var/log/fpm-php.www.log
php_admin_flag[log_errors] = on
`, poolName, version, poolName)

	poolPath := fmt.Sprintf("/etc/php/%s/fpm/pool.d/%s.conf", version, poolName)
	result := p.WriteFile(poolPath, poolConfig)
	if !result.Success {
		return result
	}
	
	// Restart PHP-FPM
	serviceName := fmt.Sprintf("php%s-fpm", version)
	return p.RestartService(serviceName)
}

// SetDefaultPHP sets the default PHP version
func (p *PHPAction) SetDefaultPHP(version string) *Result {
	// Update alternatives
	phpBinary := fmt.Sprintf("/usr/bin/php%s", version)
	if !p.FileExists(phpBinary) {
		return &Result{
			Success: false,
			Message: fmt.Sprintf("PHP %s is not installed", version),
			Error:   fmt.Errorf("PHP version not found"),
		}
	}
	
	// Update alternatives for php
	result := p.RunCommand("update-alternatives", "--install", "/usr/bin/php", "php", phpBinary, "1")
	if !result.Success {
		return result
	}
	
	// Set as default
	return p.RunCommand("update-alternatives", "--set", "php", phpBinary)
}

// Private helper methods

func (p *PHPAction) installPHPDebian(version string) *Result {
	// Add Ondrej's PHP repository for multiple versions
	if !p.FileExists("/etc/apt/sources.list.d/ondrej-ubuntu-php-*.list") {
		addRepoResult := p.RunCommand("add-apt-repository", "-y", "ppa:ondrej/php")
		if !addRepoResult.Success {
			return addRepoResult
		}
		
		updateResult := p.RunCommand("apt", "update")
		if !updateResult.Success {
			return updateResult
		}
	}
	
	// Install PHP and common extensions
	packages := []string{
		fmt.Sprintf("php%s", version),
		fmt.Sprintf("php%s-fpm", version),
		fmt.Sprintf("php%s-cli", version),
		fmt.Sprintf("php%s-common", version),
		fmt.Sprintf("php%s-mysql", version),
		fmt.Sprintf("php%s-pgsql", version),
		fmt.Sprintf("php%s-sqlite3", version),
		fmt.Sprintf("php%s-curl", version),
		fmt.Sprintf("php%s-gd", version),
		fmt.Sprintf("php%s-mbstring", version),
		fmt.Sprintf("php%s-xml", version),
		fmt.Sprintf("php%s-zip", version),
		fmt.Sprintf("php%s-bcmath", version),
		fmt.Sprintf("php%s-intl", version),
		fmt.Sprintf("php%s-json", version),
		fmt.Sprintf("php%s-opcache", version),
		fmt.Sprintf("php%s-readline", version),
	}
	
	// Filter packages that don't exist for certain versions
	if version >= "8.0" {
		// Remove json package as it's built-in in PHP 8.0+
		packages = p.filterPackages(packages, "json")
	}
	
	args := append([]string{"install", "-y"}, packages...)
	installResult := p.RunCommand("apt", args...)
	if !installResult.Success {
		return installResult
	}
	
	// Enable and start PHP-FPM
	serviceName := fmt.Sprintf("php%s-fpm", version)
	p.EnableService(serviceName)
	return p.StartService(serviceName)
}

func (p *PHPAction) installPHPRHEL(version string) *Result {
	var installCmd string
	if p.FileExists("/usr/bin/dnf") {
		installCmd = "dnf"
	} else {
		installCmd = "yum"
	}
	
	// Enable EPEL and Remi repositories
	if !p.FileExists("/etc/yum.repos.d/epel.repo") {
		epelResult := p.RunCommand(installCmd, "install", "-y", "epel-release")
		if !epelResult.Success {
			return epelResult
		}
	}
	
	if !p.FileExists("/etc/yum.repos.d/remi.repo") {
		remiResult := p.RunCommand(installCmd, "install", "-y", "https://rpms.remirepo.net/enterprise/remi-release-8.rpm")
		if !remiResult.Success {
			return remiResult
		}
	}
	
	// Enable PHP version repository
	versionNoDot := strings.Replace(version, ".", "", -1)
	enableResult := p.RunCommand(installCmd, "config-manager", "--enable", fmt.Sprintf("remi-php%s", versionNoDot))
	if !enableResult.Success {
		return enableResult
	}
	
	// Install PHP packages
	packages := []string{
		"php",
		"php-fpm",
		"php-cli",
		"php-common",
		"php-mysqlnd",
		"php-pgsql",
		"php-sqlite3",
		"php-curl",
		"php-gd",
		"php-mbstring",
		"php-xml",
		"php-zip",
		"php-bcmath",
		"php-intl",
		"php-json",
		"php-opcache",
	}
	
	args := append([]string{"install", "-y"}, packages...)
	installResult := p.RunCommand(installCmd, args...)
	if !installResult.Success {
		return installResult
	}
	
	// Enable and start PHP-FPM
	p.EnableService("php-fpm")
	return p.StartService("php-fpm")
}

func (p *PHPAction) filterPackages(packages []string, exclude string) []string {
	var filtered []string
	for _, pkg := range packages {
		if !strings.Contains(pkg, exclude) {
			filtered = append(filtered, pkg)
		}
	}
	return filtered
}