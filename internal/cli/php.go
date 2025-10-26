package cli

import (
	"easygo/pkg/actions"
	"fmt"

	"github.com/spf13/cobra"
)

var phpCmd = &cobra.Command{
	Use:   "php",
	Short: "PHP management",
	Long:  `Install, configure, and manage PHP versions and extensions.`,
}

var phpInstallCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		version := args[0]
		
		phpAction := actions.NewPHPAction()
		result := phpAction.InstallPHP(version)
		handleResult(result)
		return nil
	},
}

var phpListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed PHP versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		phpAction := actions.NewPHPAction()
		result := phpAction.GetInstalledVersions()
		
		if result.Success {
			if versions, ok := result.Data.([]*actions.PHPVersion); ok {
				fmt.Println("Installed PHP versions:")
				for _, version := range versions {
					status := "✗ Stopped"
					if version.FPMRunning {
						status = "✓ Running"
					}
					fmt.Printf("  PHP %s - FPM: %s\n", version.Version, status)
				}
			}
		} else {
			handleResult(result)
		}
		return nil
	},
}

var phpDefaultCmd = &cobra.Command{
	Use:   "default [version]",
	Short: "Set default PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		version := args[0]
		
		phpAction := actions.NewPHPAction()
		result := phpAction.SetDefaultPHP(version)
		handleResult(result)
		return nil
	},
}

var phpPoolCmd = &cobra.Command{
	Use:   "pool [version] [pool-name]",
	Short: "Configure PHP-FPM pool",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		version := args[0]
		poolName := args[1]
		
		phpAction := actions.NewPHPAction()
		result := phpAction.ConfigurePHPFPM(version, poolName)
		handleResult(result)
		return nil
	},
}

var phpAvailableCmd = &cobra.Command{
	Use:   "available",
	Short: "List available PHP versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		phpAction := actions.NewPHPAction()
		versions := phpAction.GetAvailableVersions()
		
		fmt.Println("Available PHP versions:")
		for _, version := range versions {
			fmt.Printf("  %s\n", version)
		}
		return nil
	},
}

func init() {
	phpCmd.AddCommand(phpInstallCmd)
	phpCmd.AddCommand(phpListCmd)
	phpCmd.AddCommand(phpDefaultCmd)
	phpCmd.AddCommand(phpPoolCmd)
	phpCmd.AddCommand(phpAvailableCmd)
}