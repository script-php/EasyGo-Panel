package cli

import (
	"easygo/pkg/actions"

	"github.com/spf13/cobra"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Cron job management",
	Long:  `Manage cron jobs and scheduled tasks.`,
}

var cronListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cron jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cronAction := actions.NewCronAction()
		result := cronAction.ListCronJobs()
		handleResult(result)
		return nil
	},
}

var cronSystemCmd = &cobra.Command{
	Use:   "system",
	Short: "List system cron jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cronAction := actions.NewCronAction()
		result := cronAction.ListSystemCronJobs()
		handleResult(result)
		return nil
	},
}

var cronAddCmd = &cobra.Command{
	Use:   "add [schedule] [command]",
	Short: "Add cron job",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		schedule := args[0]
		command := args[1]
		description, _ := cmd.Flags().GetString("description")
		
		cronAction := actions.NewCronAction()
		result := cronAction.AddCronJob(schedule, command, description)
		handleResult(result)
		return nil
	},
}

var cronRemoveCmd = &cobra.Command{
	Use:   "remove [command]",
	Short: "Remove cron job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		command := args[0]
		
		cronAction := actions.NewCronAction()
		result := cronAction.RemoveCronJob(command)
		handleResult(result)
		return nil
	},
}

var cronStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check cron service status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cronAction := actions.NewCronAction()
		result := cronAction.GetCronStatus()
		handleResult(result)
		return nil
	},
}

var cronMaintenanceCmd = &cobra.Command{
	Use:   "setup-maintenance",
	Short: "Setup system maintenance cron jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		
		cronAction := actions.NewCronAction()
		result := cronAction.SetupSystemMaintenance()
		handleResult(result)
		return nil
	},
}

func init() {
	cronAddCmd.Flags().StringP("description", "d", "", "Description for the cron job")
	
	cronCmd.AddCommand(cronListCmd)
	cronCmd.AddCommand(cronSystemCmd)
	cronCmd.AddCommand(cronAddCmd)
	cronCmd.AddCommand(cronRemoveCmd)
	cronCmd.AddCommand(cronStatusCmd)
	cronCmd.AddCommand(cronMaintenanceCmd)
}