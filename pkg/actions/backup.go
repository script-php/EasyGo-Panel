package actions

import (
	"fmt"
	"time"
)

// BackupAction handles backup management
type BackupAction struct {
	BaseAction
}

// NewBackupAction creates a new backup action instance
func NewBackupAction() *BackupAction {
	return &BackupAction{}
}

// BackupJob represents a backup job
type BackupJob struct {
	ID          string
	Name        string
	Type        string // files, database, full
	Source      string
	Destination string
	Schedule    string
	Enabled     bool
	LastRun     time.Time
	Status      string
}

// CreateFileBackup creates a backup of files/directories
func (b *BackupAction) CreateFileBackup(source, destination, name string) *Result {
	// Create backup directory if it doesn't exist
	if !b.DirectoryExists(destination) {
		createResult := b.CreateDirectory(destination)
		if !createResult.Success {
			return createResult
		}
	}
	
	// Create timestamped backup filename
	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("%s/%s_%s.tar.gz", destination, name, timestamp)
	
	// Create tar.gz backup
	return b.RunCommand("tar", "-czf", backupFile, "-C", source, ".")
}

// CreateDatabaseBackup creates a database backup
func (b *BackupAction) CreateDatabaseBackup(dbName, dbType, destination string) *Result {
	// Create backup directory if it doesn't exist
	if !b.DirectoryExists(destination) {
		createResult := b.CreateDirectory(destination)
		if !createResult.Success {
			return createResult
		}
	}
	
	timestamp := time.Now().Format("20060102_150405")
	
	switch dbType {
	case "mysql", "mariadb":
		backupFile := fmt.Sprintf("%s/%s_%s.sql", destination, dbName, timestamp)
		return b.RunCommand("mysqldump", dbName, ">", backupFile)
	case "postgresql":
		backupFile := fmt.Sprintf("%s/%s_%s.sql", destination, dbName, timestamp)
		return b.RunCommand("sudo", "-u", "postgres", "pg_dump", dbName, ">", backupFile)
	default:
		return &Result{
			Success: false,
			Message: fmt.Sprintf("Unsupported database type: %s", dbType),
			Error:   fmt.Errorf("unsupported database type"),
		}
	}
}

// CreateFullSystemBackup creates a full system backup
func (b *BackupAction) CreateFullSystemBackup(destination string) *Result {
	// Create backup directory if it doesn't exist
	if !b.DirectoryExists(destination) {
		createResult := b.CreateDirectory(destination)
		if !createResult.Success {
			return createResult
		}
	}
	
	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("%s/system_backup_%s.tar.gz", destination, timestamp)
	
	// Backup essential system directories
	excludeArgs := []string{
		"--exclude=/proc",
		"--exclude=/tmp",
		"--exclude=/mnt",
		"--exclude=/dev",
		"--exclude=/sys",
		"--exclude=/run",
		"--exclude=/media",
		"--exclude=/lost+found",
	}
	
	args := append([]string{"-czf", backupFile}, excludeArgs...)
	args = append(args, "/")
	
	return b.RunCommand("tar", args...)
}

// RestoreFileBackup restores files from backup
func (b *BackupAction) RestoreFileBackup(backupFile, destination string) *Result {
	if !b.FileExists(backupFile) {
		return &Result{
			Success: false,
			Message: "Backup file not found",
			Error:   fmt.Errorf("backup file does not exist"),
		}
	}
	
	// Create destination directory if it doesn't exist
	if !b.DirectoryExists(destination) {
		createResult := b.CreateDirectory(destination)
		if !createResult.Success {
			return createResult
		}
	}
	
	return b.RunCommand("tar", "-xzf", backupFile, "-C", destination)
}

// ListBackups lists available backups
func (b *BackupAction) ListBackups(backupDir string) *Result {
	if !b.DirectoryExists(backupDir) {
		return &Result{
			Success: false,
			Message: "Backup directory not found",
			Error:   fmt.Errorf("backup directory does not exist"),
		}
	}
	
	return b.RunCommand("ls", "-la", backupDir)
}

// CleanOldBackups removes old backup files
func (b *BackupAction) CleanOldBackups(backupDir string, daysToKeep int) *Result {
	if !b.DirectoryExists(backupDir) {
		return &Result{
			Success: false,
			Message: "Backup directory not found",
			Error:   fmt.Errorf("backup directory does not exist"),
		}
	}
	
	// Find and delete files older than specified days
	findArgs := []string{
		backupDir,
		"-name", "*.tar.gz",
		"-o", "-name", "*.sql",
		"-mtime", fmt.Sprintf("+%d", daysToKeep),
		"-delete",
	}
	
	return b.RunCommand("find", findArgs...)
}

// SetupAutomaticBackup sets up automatic backup via cron
func (b *BackupAction) SetupAutomaticBackup(jobName, schedule, backupType, source, destination string) *Result {
	// Create backup script
	scriptPath := fmt.Sprintf("/opt/easygo/scripts/backup_%s.sh", jobName)
	
	var script string
	switch backupType {
	case "files":
		script = fmt.Sprintf(`#!/bin/bash
# Automatic backup script for %s
LOG_FILE="/var/log/easygo/backup_%s.log"
echo "$(date): Starting backup of %s" >> $LOG_FILE
/usr/local/bin/easygo backup files %s %s >> $LOG_FILE 2>&1
echo "$(date): Backup completed" >> $LOG_FILE
`, jobName, jobName, source, source, destination)
	case "database":
		script = fmt.Sprintf(`#!/bin/bash
# Automatic database backup script for %s
LOG_FILE="/var/log/easygo/backup_%s.log"
echo "$(date): Starting database backup of %s" >> $LOG_FILE
/usr/local/bin/easygo backup database %s %s >> $LOG_FILE 2>&1
echo "$(date): Database backup completed" >> $LOG_FILE
`, jobName, jobName, source, source, destination)
	case "full":
		script = fmt.Sprintf(`#!/bin/bash
# Automatic full system backup script
LOG_FILE="/var/log/easygo/backup_full.log"
echo "$(date): Starting full system backup" >> $LOG_FILE
/usr/local/bin/easygo backup full %s >> $LOG_FILE 2>&1
echo "$(date): Full system backup completed" >> $LOG_FILE
`, destination)
	}
	
	// Create script directory
	if !b.DirectoryExists("/opt/easygo/scripts") {
		b.CreateDirectory("/opt/easygo/scripts")
	}
	
	// Write script file
	writeResult := b.WriteFile(scriptPath, script)
	if !writeResult.Success {
		return writeResult
	}
	
	// Make script executable
	chmodResult := b.RunCommand("chmod", "+x", scriptPath)
	if !chmodResult.Success {
		return chmodResult
	}
	
	// Add to crontab
	cronEntry := fmt.Sprintf("%s %s", schedule, scriptPath)
	cmd := fmt.Sprintf("(crontab -l 2>/dev/null; echo '%s') | crontab -", cronEntry)
	
	return b.RunCommand("bash", "-c", cmd)
}