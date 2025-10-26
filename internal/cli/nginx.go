package cli

import (
	"easygo/pkg/actions"

	"github.com/spf13/cobra"
)

var nginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "Nginx web server management",
	Long:  `Install, configure, and manage Nginx web server.`,
}

var nginxInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Nginx web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.InstallNginx()
		handleResult(result)
		return nil
	},
}

var nginxStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Nginx status",
	RunE: func(cmd *cobra.Command, args []string) error {
		webAction := actions.NewWebServerAction()
		result := webAction.ServiceStatus("nginx")
		handleResult(result)
		return nil
	},
}

var nginxVhostCmd = &cobra.Command{
	Use:   "vhost [domain] [document-root]",
	Short: "Create Nginx virtual host",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		domain := args[0]
		docroot := args[1]
		
		webAction := actions.NewWebServerAction()
		result := webAction.ConfigureNginxVhost(domain, docroot)
		handleResult(result)
		return nil
	},
}

var nginxStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Nginx service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.StartService("nginx")
		handleResult(result)
		return nil
	},
}

var nginxStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Nginx service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.StopService("nginx")
		handleResult(result)
		return nil
	},
}

var nginxRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart Nginx service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		webAction := actions.NewWebServerAction()
		result := webAction.RestartService("nginx")
		handleResult(result)
		return nil
	},
}

func init() {
	nginxCmd.AddCommand(nginxInstallCmd)
	nginxCmd.AddCommand(nginxStatusCmd)
	nginxCmd.AddCommand(nginxVhostCmd)
	nginxCmd.AddCommand(nginxStartCmd)
	nginxCmd.AddCommand(nginxStopCmd)
	nginxCmd.AddCommand(nginxRestartCmd)
}