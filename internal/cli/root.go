package cli

import (
	"easygo/pkg/actions"
	"easygo/pkg/auth"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easygo",
	Short: "EasyGo Panel - Web Server Management Tool",
	Long: `EasyGo Panel is a comprehensive web server management tool for Linux systems.
It provides both CLI and web interfaces for managing web servers, databases, 
mail services, DNS, SSL certificates, and more.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	
	// Add subcommands
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(apacheCmd)
	rootCmd.AddCommand(nginxCmd)
	rootCmd.AddCommand(phpCmd)
	rootCmd.AddCommand(dnsCmd)
	rootCmd.AddCommand(mailCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(sslCmd)
	rootCmd.AddCommand(firewallCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(cronCmd)
	rootCmd.AddCommand(statusCmd)
}

// requireRoot is a helper function to check root privileges
func requireRoot() error {
	return auth.RequireRoot()
}

// handleResult processes action results and displays appropriate output
func handleResult(result *actions.Result) {
	if result.Success {
		fmt.Printf("✓ %s\n", result.Message)
		if result.Data != nil {
			fmt.Printf("Data: %+v\n", result.Data)
		}
	} else {
		fmt.Printf("✗ %s\n", result.Message)
		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
		}
		os.Exit(1)
	}
}