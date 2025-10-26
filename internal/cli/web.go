package cli

import (
	"easygo/internal/web"
	"fmt"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web panel",
	Long:  `Start the EasyGo web panel interface on the specified port.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetString("port")
		host, _ := cmd.Flags().GetString("host")
		
		fmt.Printf("Starting EasyGo Web Panel on %s:%s\n", host, port)
		
		server := web.NewServer()
		return server.Start(fmt.Sprintf("%s:%s", host, port))
	},
}

func init() {
	webCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	webCmd.Flags().StringP("host", "H", "0.0.0.0", "Host to bind to")
}