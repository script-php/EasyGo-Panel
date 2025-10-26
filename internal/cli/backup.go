package cli

import (
	"easygo/pkg/actions"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup management",
	Long:  `Create, restore, and manage system backups.`,
}

var backupFilesCmd = &cobra.Command{
	Use:   "files [source] [destination] [name]",
	Short: "Create file backup",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		source := args[0]
		destination := args[1]
		name := args[2]
		
		backupAction := actions.NewBackupAction()
		result := backupAction.CreateFileBackup(source, destination, name)
		handleResult(result)
		return nil
	},
}

var backupDatabaseCmd = &cobra.Command{
	Use:   "database [db-name] [db-type] [destination]",
	Short: "Create database backup",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		dbName := args[0]
		dbType := args[1]
		destination := args[2]
		
		backupAction := actions.NewBackupAction()
		result := backupAction.CreateDatabaseBackup(dbName, dbType, destination)
		handleResult(result)
		return nil
	},
}

var backupFullCmd = &cobra.Command{
	Use:   "full [destination]",
	Short: "Create full system backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		destination := args[0]
		
		backupAction := actions.NewBackupAction()
		result := backupAction.CreateFullSystemBackup(destination)
		handleResult(result)
		return nil
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore [backup-file] [destination]",
	Short: "Restore from backup",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		backupFile := args[0]
		destination := args[1]
		
		backupAction := actions.NewBackupAction()
		result := backupAction.RestoreFileBackup(backupFile, destination)
		handleResult(result)
		return nil
	},
}

var backupListCmd = &cobra.Command{
	Use:   "list [backup-dir]",
	Short: "List available backups",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupDir := args[0]
		
		backupAction := actions.NewBackupAction()
		result := backupAction.ListBackups(backupDir)
		handleResult(result)
		return nil
	},
}

var backupCleanCmd = &cobra.Command{
	Use:   "clean [backup-dir] [days-to-keep]",
	Short: "Clean old backups",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		backupDir := args[0]
		// Parse days to keep - simplified for now
		daysToKeep := 30 // default
		
		backupAction := actions.NewBackupAction()
		result := backupAction.CleanOldBackups(backupDir, daysToKeep)
		handleResult(result)
		return nil
	},
}

func init() {
	backupCmd.AddCommand(backupFilesCmd)
	backupCmd.AddCommand(backupDatabaseCmd)
	backupCmd.AddCommand(backupFullCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupCleanCmd)
}