package actions

import (
	"fmt"
	"strings"
	"time"
)

// CronAction handles cron job management
type CronAction struct {
	BaseAction
}

// NewCronAction creates a new cron action instance
func NewCronAction() *CronAction {
	return &CronAction{}
}

// CronJob represents a cron job
type CronJob struct {
	ID          string
	Schedule    string
	Command     string
	User        string
	Description string
	Enabled     bool
	LastRun     time.Time
	NextRun     time.Time
}

// ListCronJobs lists all cron jobs for the current user
func (c *CronAction) ListCronJobs() *Result {
	return c.RunCommand("crontab", "-l")
}

// ListSystemCronJobs lists system-wide cron jobs
func (c *CronAction) ListSystemCronJobs() *Result {
	result := c.RunCommand("cat", "/etc/crontab")
	if !result.Success {
		return result
	}
	
	// Also check /etc/cron.d/ directory
	cronDResult := c.RunCommand("ls", "-la", "/etc/cron.d/")
	if cronDResult.Success {
		result.Message += "\n\n--- /etc/cron.d/ ---\n" + cronDResult.Message
	}
	
	return result
}

// AddCronJob adds a new cron job
func (c *CronAction) AddCronJob(schedule, command, description string) *Result {
	// Validate cron schedule format (basic validation)
	if !c.isValidCronSchedule(schedule) {
		return &Result{
			Success: false,
			Message: "Invalid cron schedule format",
			Error:   fmt.Errorf("invalid schedule format"),
		}
	}
	
	// Add comment with description if provided
	var cronEntry string
	if description != "" {
		cronEntry = fmt.Sprintf("# %s\n%s %s", description, schedule, command)
	} else {
		cronEntry = fmt.Sprintf("%s %s", schedule, command)
	}
	
	// Add to crontab
	cmd := fmt.Sprintf("(crontab -l 2>/dev/null; echo '%s') | crontab -", cronEntry)
	return c.RunCommand("bash", "-c", cmd)
}

// RemoveCronJob removes a cron job by matching the command
func (c *CronAction) RemoveCronJob(command string) *Result {
	// Get current crontab
	currentResult := c.RunCommand("crontab", "-l")
	if !currentResult.Success {
		return &Result{
			Success: false,
			Message: "No crontab found",
			Error:   fmt.Errorf("no crontab exists"),
		}
	}
	
	// Filter out the matching command
	lines := strings.Split(currentResult.Message, "\n")
	var filteredLines []string
	
	for _, line := range lines {
		if !strings.Contains(line, command) {
			filteredLines = append(filteredLines, line)
		}
	}
	
	// Write back the filtered crontab
	newCrontab := strings.Join(filteredLines, "\n")
	cmd := fmt.Sprintf("echo '%s' | crontab -", newCrontab)
	
	return c.RunCommand("bash", "-c", cmd)
}

// EnableCronService enables the cron service
func (c *CronAction) EnableCronService() *Result {
	// Enable cron service (varies by system)
	enableResult := c.EnableService("cron")
	if !enableResult.Success {
		// Try alternative service names
		enableResult = c.EnableService("crond")
	}
	
	return enableResult
}

// StartCronService starts the cron service
func (c *CronAction) StartCronService() *Result {
	startResult := c.StartService("cron")
	if !startResult.Success {
		// Try alternative service names
		startResult = c.StartService("crond")
	}
	
	return startResult
}

// GetCronStatus gets the status of the cron service
func (c *CronAction) GetCronStatus() *Result {
	statusResult := c.ServiceStatus("cron")
	if !statusResult.Success {
		// Try alternative service names
		statusResult = c.ServiceStatus("crond")
	}
	
	return statusResult
}

// AddSystemCronJob adds a system-wide cron job
func (c *CronAction) AddSystemCronJob(name, schedule, user, command, description string) *Result {
	cronFile := fmt.Sprintf("/etc/cron.d/%s", name)
	
	cronContent := fmt.Sprintf(`# %s
%s %s %s
`, description, schedule, user, command)
	
	return c.WriteFile(cronFile, cronContent)
}

// RemoveSystemCronJob removes a system-wide cron job
func (c *CronAction) RemoveSystemCronJob(name string) *Result {
	cronFile := fmt.Sprintf("/etc/cron.d/%s", name)
	return c.RunCommand("rm", "-f", cronFile)
}

// AddDailyCronJob adds a job to run daily
func (c *CronAction) AddDailyCronJob(hour, minute int, command, description string) *Result {
	schedule := fmt.Sprintf("%d %d * * *", minute, hour)
	return c.AddCronJob(schedule, command, description)
}

// AddWeeklyCronJob adds a job to run weekly
func (c *CronAction) AddWeeklyCronJob(weekday, hour, minute int, command, description string) *Result {
	schedule := fmt.Sprintf("%d %d * * %d", minute, hour, weekday)
	return c.AddCronJob(schedule, command, description)
}

// AddMonthlyCronJob adds a job to run monthly
func (c *CronAction) AddMonthlyCronJob(day, hour, minute int, command, description string) *Result {
	schedule := fmt.Sprintf("%d %d %d * *", minute, hour, day)
	return c.AddCronJob(schedule, command, description)
}

// SetupLogRotation sets up log rotation for EasyGo Panel
func (c *CronAction) SetupLogRotation() *Result {
	logRotateConfig := `/var/log/easygo/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 root root
    postrotate
        systemctl reload easygo 2>/dev/null || true
    endscript
}`

	return c.WriteFile("/etc/logrotate.d/easygo", logRotateConfig)
}

// SetupSystemMaintenance sets up basic system maintenance cron jobs
func (c *CronAction) SetupSystemMaintenance() *Result {
	// Add system update check (weekly)
	updateCheckResult := c.AddWeeklyCronJob(0, 2, 0, "/usr/bin/apt update && /usr/bin/apt list --upgradable", "Weekly system update check")
	if !updateCheckResult.Success {
		return updateCheckResult
	}
	
	// Add log cleanup (daily)
	logCleanupResult := c.AddDailyCronJob(1, 0, "find /var/log -name '*.log' -mtime +30 -delete", "Daily old log cleanup")
	if !logCleanupResult.Success {
		return logCleanupResult
	}
	
	// Add temporary file cleanup (daily)
	tempCleanupResult := c.AddDailyCronJob(1, 30, "find /tmp -type f -mtime +7 -delete", "Daily temporary file cleanup")
	if !tempCleanupResult.Success {
		return tempCleanupResult
	}
	
	return &Result{
		Success: true,
		Message: "System maintenance cron jobs added successfully",
	}
}

// Private helper methods

func (c *CronAction) isValidCronSchedule(schedule string) bool {
	// Basic validation - check if it has 5 parts (minute hour day month weekday)
	parts := strings.Fields(schedule)
	return len(parts) == 5
}