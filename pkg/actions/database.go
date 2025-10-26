package actions

import (
	"fmt"
	"time"
)

// DatabaseAction handles database management
type DatabaseAction struct {
	BaseAction
}

// NewDatabaseAction creates a new database action instance
func NewDatabaseAction() *DatabaseAction {
	return &DatabaseAction{}
}

// Database represents a database
type Database struct {
	Name     string
	Type     string // mysql, mariadb, postgresql
	Size     string
	User     string
	Created  time.Time
	Status   string
}

// InstallMariaDB installs MariaDB server
func (d *DatabaseAction) InstallMariaDB() *Result {
	if d.FileExists("/usr/bin/apt") {
		return d.installMariaDBDebian()
	} else if d.FileExists("/usr/bin/yum") || d.FileExists("/usr/bin/dnf") {
		return d.installMariaDBRHEL()
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported Linux distribution",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}

// InstallPostgreSQL installs PostgreSQL server
func (d *DatabaseAction) InstallPostgreSQL() *Result {
	if d.FileExists("/usr/bin/apt") {
		return d.installPostgreSQLDebian()
	} else if d.FileExists("/usr/bin/yum") || d.FileExists("/usr/bin/dnf") {
		return d.installPostgreSQLRHEL()
	}
	
	return &Result{
		Success: false,
		Message: "Unsupported Linux distribution",
		Error:   fmt.Errorf("unsupported package manager"),
	}
}

// CreateDatabase creates a new database
func (d *DatabaseAction) CreateDatabase(name, dbType, username, password string) *Result {
	switch dbType {
	case "mysql", "mariadb":
		return d.createMySQLDatabase(name, username, password)
	case "postgresql":
		return d.createPostgreSQLDatabase(name, username, password)
	default:
		return &Result{
			Success: false,
			Message: fmt.Sprintf("Unsupported database type: %s", dbType),
			Error:   fmt.Errorf("unsupported database type"),
		}
	}
}

// DropDatabase drops a database
func (d *DatabaseAction) DropDatabase(name, dbType string) *Result {
	switch dbType {
	case "mysql", "mariadb":
		return d.dropMySQLDatabase(name)
	case "postgresql":
		return d.dropPostgreSQLDatabase(name)
	default:
		return &Result{
			Success: false,
			Message: fmt.Sprintf("Unsupported database type: %s", dbType),
			Error:   fmt.Errorf("unsupported database type"),
		}
	}
}

// ListDatabases lists all databases
func (d *DatabaseAction) ListDatabases(dbType string) *Result {
	switch dbType {
	case "mysql", "mariadb":
		return d.listMySQLDatabases()
	case "postgresql":
		return d.listPostgreSQLDatabases()
	default:
		return &Result{
			Success: false,
			Message: fmt.Sprintf("Unsupported database type: %s", dbType),
			Error:   fmt.Errorf("unsupported database type"),
		}
	}
}

// BackupDatabase creates a backup of a database
func (d *DatabaseAction) BackupDatabase(name, dbType, backupPath string) *Result {
	switch dbType {
	case "mysql", "mariadb":
		return d.backupMySQLDatabase(name, backupPath)
	case "postgresql":
		return d.backupPostgreSQLDatabase(name, backupPath)
	default:
		return &Result{
			Success: false,
			Message: fmt.Sprintf("Unsupported database type: %s", dbType),
			Error:   fmt.Errorf("unsupported database type"),
		}
	}
}

// RestoreDatabase restores a database from backup
func (d *DatabaseAction) RestoreDatabase(name, dbType, backupPath string) *Result {
	switch dbType {
	case "mysql", "mariadb":
		return d.restoreMySQLDatabase(name, backupPath)
	case "postgresql":
		return d.restorePostgreSQLDatabase(name, backupPath)
	default:
		return &Result{
			Success: false,
			Message: fmt.Sprintf("Unsupported database type: %s", dbType),
			Error:   fmt.Errorf("unsupported database type"),
		}
	}
}

// InstallPHPMyAdmin installs phpMyAdmin
func (d *DatabaseAction) InstallPHPMyAdmin() *Result {
	if d.FileExists("/usr/bin/apt") {
		// Pre-configure phpMyAdmin for automatic installation
		d.RunCommand("echo", "phpmyadmin", "phpmyadmin/dbconfig-install", "boolean", "true", "|", "debconf-set-selections")
		d.RunCommand("echo", "phpmyadmin", "phpmyadmin/app-password-confirm", "password", "admin", "|", "debconf-set-selections")
		d.RunCommand("echo", "phpmyadmin", "phpmyadmin/mysql/admin-pass", "password", "", "|", "debconf-set-selections")
		d.RunCommand("echo", "phpmyadmin", "phpmyadmin/reconfigure-webserver", "multiselect", "apache2", "|", "debconf-set-selections")
		
		return d.RunCommand("apt", "install", "-y", "phpmyadmin", "php-mbstring", "php-zip", "php-gd", "php-json", "php-curl")
	}
	
	return &Result{
		Success: false,
		Message: "Manual phpMyAdmin installation required for this distribution",
		Error:   fmt.Errorf("automatic installation not supported"),
	}
}

// Private helper methods for MariaDB/MySQL

func (d *DatabaseAction) installMariaDBDebian() *Result {
	updateResult := d.RunCommand("apt", "update")
	if !updateResult.Success {
		return updateResult
	}
	
	installResult := d.RunCommand("apt", "install", "-y", "mariadb-server", "mariadb-client")
	if !installResult.Success {
		return installResult
	}
	
	// Start and enable MariaDB
	d.EnableService("mariadb")
	startResult := d.StartService("mariadb")
	if !startResult.Success {
		return startResult
	}
	
	// Secure installation
	return d.secureMariaDBInstallation()
}

func (d *DatabaseAction) installMariaDBRHEL() *Result {
	var installCmd string
	if d.FileExists("/usr/bin/dnf") {
		installCmd = "dnf"
	} else {
		installCmd = "yum"
	}
	
	installResult := d.RunCommand(installCmd, "install", "-y", "mariadb-server", "mariadb")
	if !installResult.Success {
		return installResult
	}
	
	// Start and enable MariaDB
	d.EnableService("mariadb")
	startResult := d.StartService("mariadb")
	if !startResult.Success {
		return startResult
	}
	
	return d.secureMariaDBInstallation()
}

func (d *DatabaseAction) secureMariaDBInstallation() *Result {
	// This is a simplified version - in production, you'd want proper security setup
	return d.RunCommand("mysql", "-e", "UPDATE mysql.user SET Password=PASSWORD('rootpassword') WHERE User='root'; FLUSH PRIVILEGES;")
}

func (d *DatabaseAction) createMySQLDatabase(name, username, password string) *Result {
	// Create database
	createDBResult := d.RunCommand("mysql", "-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", name))
	if !createDBResult.Success {
		return createDBResult
	}
	
	// Create user and grant privileges
	createUserSQL := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'localhost' IDENTIFIED BY '%s'; GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost'; FLUSH PRIVILEGES;", username, password, name, username)
	return d.RunCommand("mysql", "-e", createUserSQL)
}

func (d *DatabaseAction) dropMySQLDatabase(name string) *Result {
	return d.RunCommand("mysql", "-e", fmt.Sprintf("DROP DATABASE IF EXISTS %s;", name))
}

func (d *DatabaseAction) listMySQLDatabases() *Result {
	return d.RunCommand("mysql", "-e", "SHOW DATABASES;")
}

func (d *DatabaseAction) backupMySQLDatabase(name, backupPath string) *Result {
	return d.RunCommand("mysqldump", name, ">", backupPath)
}

func (d *DatabaseAction) restoreMySQLDatabase(name, backupPath string) *Result {
	return d.RunCommand("mysql", name, "<", backupPath)
}

// Private helper methods for PostgreSQL

func (d *DatabaseAction) installPostgreSQLDebian() *Result {
	updateResult := d.RunCommand("apt", "update")
	if !updateResult.Success {
		return updateResult
	}
	
	installResult := d.RunCommand("apt", "install", "-y", "postgresql", "postgresql-contrib")
	if !installResult.Success {
		return installResult
	}
	
	// Start and enable PostgreSQL
	d.EnableService("postgresql")
	return d.StartService("postgresql")
}

func (d *DatabaseAction) installPostgreSQLRHEL() *Result {
	var installCmd string
	if d.FileExists("/usr/bin/dnf") {
		installCmd = "dnf"
	} else {
		installCmd = "yum"
	}
	
	installResult := d.RunCommand(installCmd, "install", "-y", "postgresql-server", "postgresql-contrib")
	if !installResult.Success {
		return installResult
	}
	
	// Initialize database
	initResult := d.RunCommand("postgresql-setup", "initdb")
	if !initResult.Success {
		return initResult
	}
	
	// Start and enable PostgreSQL
	d.EnableService("postgresql")
	return d.StartService("postgresql")
}

func (d *DatabaseAction) createPostgreSQLDatabase(name, username, password string) *Result {
	// Create user
	createUserResult := d.RunCommand("sudo", "-u", "postgres", "createuser", username)
	if !createUserResult.Success {
		return createUserResult
	}
	
	// Create database
	createDBResult := d.RunCommand("sudo", "-u", "postgres", "createdb", "-O", username, name)
	if !createDBResult.Success {
		return createDBResult
	}
	
	// Set password
	setPasswordSQL := fmt.Sprintf("ALTER USER %s PASSWORD '%s';", username, password)
	return d.RunCommand("sudo", "-u", "postgres", "psql", "-c", setPasswordSQL)
}

func (d *DatabaseAction) dropPostgreSQLDatabase(name string) *Result {
	return d.RunCommand("sudo", "-u", "postgres", "dropdb", name)
}

func (d *DatabaseAction) listPostgreSQLDatabases() *Result {
	return d.RunCommand("sudo", "-u", "postgres", "psql", "-l")
}

func (d *DatabaseAction) backupPostgreSQLDatabase(name, backupPath string) *Result {
	return d.RunCommand("sudo", "-u", "postgres", "pg_dump", name, ">", backupPath)
}

func (d *DatabaseAction) restorePostgreSQLDatabase(name, backupPath string) *Result {
	return d.RunCommand("sudo", "-u", "postgres", "psql", name, "<", backupPath)
}