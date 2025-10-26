package cli

import (
	"easygo/pkg/actions"

	"github.com/spf13/cobra"
)

var apacheCmd = &cobra.Command{
	Use:   "apache",
	Short: "Apache web server management",
	Long:  `Install, configure, and manage Apache web server.`,
}

var apacheInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Apache web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.InstallApache()
		handleResult(result)
		return nil
	},
}

var apacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Apache status",
	RunE: func(cmd *cobra.Command, args []string) error {
		webAction := actions.NewWebServerAction()
		result := webAction.ServiceStatus("apache2")
		handleResult(result)
		return nil
	},
}

var apacheVhostCmd = &cobra.Command{
	Use:   "vhost [domain] [document-root]",
	Short: "Create Apache virtual host",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		domain := args[0]
		docroot := args[1]
		
		webAction := actions.NewWebServerAction()
		result := webAction.ConfigureApacheVhost(domain, docroot)
		handleResult(result)
		return nil
	},
}

var apacheStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Apache service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.StartService("apache2")
		handleResult(result)
		return nil
	},
}

var apacheStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Apache service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.StopService("apache2")
		handleResult(result)
		return nil
	},
}

var apacheRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart Apache service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.RestartService("apache2")
		handleResult(result)
		return nil
	},
}

func init() {
	apacheCmd.AddCommand(apacheInstallCmd)
	apacheCmd.AddCommand(apacheStatusCmd)
	apacheCmd.AddCommand(apacheVhostCmd)
	apacheCmd.AddCommand(apacheStartCmd)
	apacheCmd.AddCommand(apacheStopCmd)
	apacheCmd.AddCommand(apacheRestartCmd)
}