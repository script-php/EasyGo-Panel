package actions

import (
	"fmt"
	"os/exec"
	"strings"
)

// Result represents the result of an action
type Result struct {
	Success bool
	Message string
	Error   error
	Data    interface{}
}

// Service represents a system service
type Service struct {
	Name    string
	Status  string
	Enabled bool
}

// BaseAction provides common functionality for all actions
type BaseAction struct{}

// RunCommand executes a system command
func (ba *BaseAction) RunCommand(command string, args ...string) *Result {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return &Result{
			Success: false,
			Message: string(output),
			Error:   err,
		}
	}
	
	return &Result{
		Success: true,
		Message: string(output),
		Error:   nil,
	}
}

// ServiceStatus checks the status of a service
func (ba *BaseAction) ServiceStatus(serviceName string) *Result {
	result := ba.RunCommand("systemctl", "is-active", serviceName)
	if result.Error != nil {
		return result
	}
	
	enabledResult := ba.RunCommand("systemctl", "is-enabled", serviceName)
	enabled := enabledResult.Success && strings.TrimSpace(enabledResult.Message) == "enabled"
	
	service := &Service{
		Name:    serviceName,
		Status:  strings.TrimSpace(result.Message),
		Enabled: enabled,
	}
	
	return &Result{
		Success: true,
		Message: fmt.Sprintf("Service %s is %s", serviceName, service.Status),
		Data:    service,
	}
}

// StartService starts a system service
func (ba *BaseAction) StartService(serviceName string) *Result {
	return ba.RunCommand("systemctl", "start", serviceName)
}

// StopService stops a system service
func (ba *BaseAction) StopService(serviceName string) *Result {
	return ba.RunCommand("systemctl", "stop", serviceName)
}

// EnableService enables a system service
func (ba *BaseAction) EnableService(serviceName string) *Result {
	return ba.RunCommand("systemctl", "enable", serviceName)
}

// DisableService disables a system service
func (ba *BaseAction) DisableService(serviceName string) *Result {
	return ba.RunCommand("systemctl", "disable", serviceName)
}

// RestartService restarts a system service
func (ba *BaseAction) RestartService(serviceName string) *Result {
	return ba.RunCommand("systemctl", "restart", serviceName)
}

// ReloadService reloads a system service
func (ba *BaseAction) ReloadService(serviceName string) *Result {
	return ba.RunCommand("systemctl", "reload", serviceName)
}

// FileExists checks if a file exists
func (ba *BaseAction) FileExists(path string) bool {
	result := ba.RunCommand("test", "-f", path)
	return result.Success
}

// DirectoryExists checks if a directory exists
func (ba *BaseAction) DirectoryExists(path string) bool {
	result := ba.RunCommand("test", "-d", path)
	return result.Success
}

// CreateDirectory creates a directory
func (ba *BaseAction) CreateDirectory(path string) *Result {
	return ba.RunCommand("mkdir", "-p", path)
}

// WriteFile writes content to a file
func (ba *BaseAction) WriteFile(path, content string) *Result {
	cmd := exec.Command("tee", path)
	cmd.Stdin = strings.NewReader(content)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return &Result{
			Success: false,
			Message: string(output),
			Error:   err,
		}
	}
	
	return &Result{
		Success: true,
		Message: "File written successfully",
	}
}