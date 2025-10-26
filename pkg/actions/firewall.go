package actions

import (
	"fmt"
)

// FirewallAction handles firewall and security management
type FirewallAction struct {
	BaseAction
}

// NewFirewallAction creates a new firewall action instance
func NewFirewallAction() *FirewallAction {
	return &FirewallAction{}
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	ID       string
	Protocol string
	Port     string
	Source   string
	Target   string
	Action   string
}

// InstallFirewall installs and configures basic firewall
func (f *FirewallAction) InstallFirewall() *Result {
	// Install iptables-persistent for Debian/Ubuntu
	if f.FileExists("/usr/bin/apt") {
		installResult := f.RunCommand("apt", "install", "-y", "iptables-persistent")
		if !installResult.Success {
			return installResult
		}
	}
	
	// Set up basic firewall rules
	return f.SetupBasicRules()
}

// SetupBasicRules configures basic firewall rules
func (f *FirewallAction) SetupBasicRules() *Result {
	rules := [][]string{
		// Allow loopback
		{"iptables", "-A", "INPUT", "-i", "lo", "-j", "ACCEPT"},
		// Allow established connections
		{"iptables", "-A", "INPUT", "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT"},
		// Allow SSH
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "22", "-j", "ACCEPT"},
		// Allow HTTP
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "80", "-j", "ACCEPT"},
		// Allow HTTPS
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "443", "-j", "ACCEPT"},
		// Allow EasyGo Panel
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "8080", "-j", "ACCEPT"},
		// Drop all other input
		{"iptables", "-P", "INPUT", "DROP"},
		// Allow all output
		{"iptables", "-P", "OUTPUT", "ACCEPT"},
		// Allow all forward
		{"iptables", "-P", "FORWARD", "ACCEPT"},
	}
	
	for _, rule := range rules {
		result := f.RunCommand(rule[0], rule[1:]...)
		if !result.Success {
			return result
		}
	}
	
	// Save rules
	return f.SaveRules()
}

// AddRule adds a new firewall rule
func (f *FirewallAction) AddRule(protocol, port, source, action string) *Result {
	var args []string
	
	if action == "allow" {
		args = []string{"-A", "INPUT", "-p", protocol}
		if port != "" {
			args = append(args, "--dport", port)
		}
		if source != "" {
			args = append(args, "-s", source)
		}
		args = append(args, "-j", "ACCEPT")
	} else if action == "deny" {
		args = []string{"-A", "INPUT", "-p", protocol}
		if port != "" {
			args = append(args, "--dport", port)
		}
		if source != "" {
			args = append(args, "-s", source)
		}
		args = append(args, "-j", "DROP")
	}
	
	result := f.RunCommand("iptables", args...)
	if !result.Success {
		return result
	}
	
	return f.SaveRules()
}

// RemoveRule removes a firewall rule
func (f *FirewallAction) RemoveRule(ruleSpec string) *Result {
	// Parse rule specification and remove it
	// This is a simplified implementation
	result := f.RunCommand("iptables", "-D", "INPUT", ruleSpec)
	if !result.Success {
		return result
	}
	
	return f.SaveRules()
}

// ListRules lists current firewall rules
func (f *FirewallAction) ListRules() *Result {
	return f.RunCommand("iptables", "-L", "-n", "--line-numbers")
}

// InstallFail2Ban installs and configures Fail2ban
func (f *FirewallAction) InstallFail2Ban() *Result {
	var installResult *Result
	
	if f.FileExists("/usr/bin/apt") {
		installResult = f.RunCommand("apt", "install", "-y", "fail2ban")
	} else if f.FileExists("/usr/bin/yum") || f.FileExists("/usr/bin/dnf") {
		var installCmd string
		if f.FileExists("/usr/bin/dnf") {
			installCmd = "dnf"
		} else {
			installCmd = "yum"
		}
		installResult = f.RunCommand(installCmd, "install", "-y", "fail2ban")
	} else {
		return &Result{
			Success: false,
			Message: "Unsupported package manager",
			Error:   fmt.Errorf("unsupported package manager"),
		}
	}
	
	if !installResult.Success {
		return installResult
	}
	
	// Configure Fail2ban
	return f.ConfigureFail2Ban()
}

// ConfigureFail2Ban configures Fail2ban with basic jails
func (f *FirewallAction) ConfigureFail2Ban() *Result {
	config := `[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5
backend = auto

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3

[apache-auth]
enabled = true
port = http,https
filter = apache-auth
logpath = /var/log/apache2/*error.log
maxretry = 6

[apache-badbots]
enabled = true
port = http,https
filter = apache-badbots
logpath = /var/log/apache2/*access.log
maxretry = 2

[nginx-http-auth]
enabled = true
port = http,https
filter = nginx-http-auth
logpath = /var/log/nginx/error.log
maxretry = 6

[nginx-badbots]
enabled = true
port = http,https
filter = nginx-badbots
logpath = /var/log/nginx/*access.log
maxretry = 2`

	result := f.WriteFile("/etc/fail2ban/jail.local", config)
	if !result.Success {
		return result
	}
	
	// Start and enable Fail2ban
	f.EnableService("fail2ban")
	return f.StartService("fail2ban")
}

// GetFail2BanStatus gets Fail2ban status
func (f *FirewallAction) GetFail2BanStatus() *Result {
	return f.RunCommand("fail2ban-client", "status")
}

// UnbanIP unbans an IP address from Fail2ban
func (f *FirewallAction) UnbanIP(ip, jail string) *Result {
	return f.RunCommand("fail2ban-client", "set", jail, "unbanip", ip)
}

// BanIP bans an IP address
func (f *FirewallAction) BanIP(ip, jail string) *Result {
	return f.RunCommand("fail2ban-client", "set", jail, "banip", ip)
}

// InstallIPSet installs and configures IPSet for IP lists
func (f *FirewallAction) InstallIPSet() *Result {
	if f.FileExists("/usr/bin/apt") {
		return f.RunCommand("apt", "install", "-y", "ipset")
	} else if f.FileExists("/usr/bin/yum") || f.FileExists("/usr/bin/dnf") {
		var installCmd string
		if f.FileExists("/usr/bin/dnf") {
			installCmd = "dnf"
		} else {
			installCmd = "yum"
		}
		return f.RunCommand(installCmd, "install", "-y", "ipset")
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported package manager",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}

// CreateIPSet creates a new IP set
func (f *FirewallAction) CreateIPSet(name, setType string) *Result {
	return f.RunCommand("ipset", "create", name, setType)
}

// AddToIPSet adds an IP to an IP set
func (f *FirewallAction) AddToIPSet(setName, ip string) *Result {
	return f.RunCommand("ipset", "add", setName, ip)
}

// RemoveFromIPSet removes an IP from an IP set
func (f *FirewallAction) RemoveFromIPSet(setName, ip string) *Result {
	return f.RunCommand("ipset", "del", setName, ip)
}

// ListIPSets lists all IP sets
func (f *FirewallAction) ListIPSets() *Result {
	return f.RunCommand("ipset", "list")
}

// SaveRules saves current iptables rules
func (f *FirewallAction) SaveRules() *Result {
	if f.FileExists("/usr/sbin/iptables-save") {
		return f.RunCommand("sh", "-c", "iptables-save > /etc/iptables/rules.v4")
	}
	return &Result{
		Success: true,
		Message: "Rules saved in memory",
	}
}

// RestoreRules restores iptables rules from file
func (f *FirewallAction) RestoreRules() *Result {
	if f.FileExists("/etc/iptables/rules.v4") {
		return f.RunCommand("iptables-restore", "/etc/iptables/rules.v4")
	}
	return &Result{
		Success: false,
		Message: "No saved rules found",
		Error:   fmt.Errorf("rules file not found"),
	}
}

// FlushRules flushes all iptables rules
func (f *FirewallAction) FlushRules() *Result {
	// Set default policies to ACCEPT before flushing
	f.RunCommand("iptables", "-P", "INPUT", "ACCEPT")
	f.RunCommand("iptables", "-P", "OUTPUT", "ACCEPT")
	f.RunCommand("iptables", "-P", "FORWARD", "ACCEPT")
	
	// Flush all chains
	result := f.RunCommand("iptables", "-F")
	if !result.Success {
		return result
	}
	
	// Delete all custom chains
	return f.RunCommand("iptables", "-X")
}