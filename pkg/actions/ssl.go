package actions

import (
	"fmt"
	"strings"
)

// SSLAction handles SSL certificate management
type SSLAction struct {
	BaseAction
}

// NewSSLAction creates a new SSL action instance
func NewSSLAction() *SSLAction {
	return &SSLAction{}
}

// Certificate represents an SSL certificate
type Certificate struct {
	Domain     string
	Type       string // single, wildcard, multi-domain
	Issuer     string
	ValidFrom  string
	ValidUntil string
	Status     string
	AutoRenew  bool
}

// InstallCertbot installs certbot for Let's Encrypt
func (s *SSLAction) InstallCertbot() *Result {
	if s.FileExists("/usr/bin/apt") {
		return s.installCertbotDebian()
	} else if s.FileExists("/usr/bin/yum") || s.FileExists("/usr/bin/dnf") {
		return s.installCertbotRHEL()
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported Linux distribution",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}

// IssueCertificate issues a Let's Encrypt certificate
func (s *SSLAction) IssueCertificate(domain, email, webroot string) *Result {
	// Check if certbot is installed
	if !s.FileExists("/usr/bin/certbot") {
		installResult := s.InstallCertbot()
		if !installResult.Success {
			return installResult
		}
	}
	
	// Issue certificate using webroot method
	args := []string{
		"certonly",
		"--webroot",
		"-w", webroot,
		"-d", domain,
		"--email", email,
		"--agree-tos",
		"--non-interactive",
	}
	
	return s.RunCommand("certbot", args...)
}

// IssueWildcardCertificate issues a wildcard certificate using DNS challenge
func (s *SSLAction) IssueWildcardCertificate(domain, email, dnsProvider string) *Result {
	if !s.FileExists("/usr/bin/certbot") {
		installResult := s.InstallCertbot()
		if !installResult.Success {
			return installResult
		}
	}
	
	// Install DNS plugin if needed
	pluginResult := s.installDNSPlugin(dnsProvider)
	if !pluginResult.Success {
		return pluginResult
	}
	
	args := []string{
		"certonly",
		"--dns-" + dnsProvider,
		"-d", "*." + domain,
		"-d", domain,
		"--email", email,
		"--agree-tos",
		"--non-interactive",
	}
	
	return s.RunCommand("certbot", args...)
}

// RenewCertificates renews all certificates
func (s *SSLAction) RenewCertificates() *Result {
	return s.RunCommand("certbot", "renew", "--quiet")
}

// ListCertificates lists all certificates
func (s *SSLAction) ListCertificates() *Result {
	result := s.RunCommand("certbot", "certificates")
	if !result.Success {
		return result
	}
	
	// Parse certbot output to extract certificate information
	var certificates []*Certificate
	lines := strings.Split(result.Message, "\n")
	
	var currentCert *Certificate
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Certificate Name:") {
			if currentCert != nil {
				certificates = append(certificates, currentCert)
			}
			currentCert = &Certificate{
				Domain: strings.TrimSpace(strings.TrimPrefix(line, "Certificate Name:")),
				Issuer: "Let's Encrypt",
			}
		} else if currentCert != nil {
			if strings.HasPrefix(line, "Domains:") {
				domains := strings.TrimSpace(strings.TrimPrefix(line, "Domains:"))
				if strings.Contains(domains, "*") {
					currentCert.Type = "wildcard"
				} else if strings.Contains(domains, " ") {
					currentCert.Type = "multi-domain"
				} else {
					currentCert.Type = "single"
				}
			} else if strings.HasPrefix(line, "Expiry Date:") {
				currentCert.ValidUntil = strings.TrimSpace(strings.TrimPrefix(line, "Expiry Date:"))
				currentCert.Status = "valid"
			}
		}
	}
	
	if currentCert != nil {
		certificates = append(certificates, currentCert)
	}
	
	return &Result{
		Success: true,
		Message: fmt.Sprintf("Found %d certificates", len(certificates)),
		Data:    certificates,
	}
}

// RevokeCertificate revokes a certificate
func (s *SSLAction) RevokeCertificate(domain string) *Result {
	return s.RunCommand("certbot", "revoke", "--cert-name", domain)
}

// SetupAutoRenewal sets up automatic certificate renewal
func (s *SSLAction) SetupAutoRenewal() *Result {
	// Create cron job for automatic renewal
	cronEntry := "0 12 * * * /usr/bin/certbot renew --quiet"
	
	// Check if cron entry already exists
	checkResult := s.RunCommand("crontab", "-l")
	if checkResult.Success && strings.Contains(checkResult.Message, "certbot renew") {
		return &Result{
			Success: true,
			Message: "Auto-renewal is already configured",
		}
	}
	
	// Add to crontab
	cmd := fmt.Sprintf("(crontab -l 2>/dev/null; echo '%s') | crontab -", cronEntry)
	return s.RunCommand("bash", "-c", cmd)
}

// Private helper methods

func (s *SSLAction) installCertbotDebian() *Result {
	updateResult := s.RunCommand("apt", "update")
	if !updateResult.Success {
		return updateResult
	}
	
	return s.RunCommand("apt", "install", "-y", "certbot")
}

func (s *SSLAction) installCertbotRHEL() *Result {
	var installCmd string
	if s.FileExists("/usr/bin/dnf") {
		installCmd = "dnf"
	} else {
		installCmd = "yum"
	}
	
	// Install EPEL repository first
	epelResult := s.RunCommand(installCmd, "install", "-y", "epel-release")
	if !epelResult.Success {
		return epelResult
	}
	
	return s.RunCommand(installCmd, "install", "-y", "certbot")
}

func (s *SSLAction) installDNSPlugin(provider string) *Result {
	var packageName string
	
	switch provider {
	case "cloudflare":
		packageName = "python3-certbot-dns-cloudflare"
	case "route53":
		packageName = "python3-certbot-dns-route53"
	case "digitalocean":
		packageName = "python3-certbot-dns-digitalocean"
	default:
		return &Result{
			Success: false,
			Message: fmt.Sprintf("Unsupported DNS provider: %s", provider),
			Error:   fmt.Errorf("unsupported DNS provider"),
		}
	}
	
	if s.FileExists("/usr/bin/apt") {
		return s.RunCommand("apt", "install", "-y", packageName)
	} else if s.FileExists("/usr/bin/yum") || s.FileExists("/usr/bin/dnf") {
		var installCmd string
		if s.FileExists("/usr/bin/dnf") {
			installCmd = "dnf"
		} else {
			installCmd = "yum"
		}
		return s.RunCommand(installCmd, "install", "-y", packageName)
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported package manager",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}