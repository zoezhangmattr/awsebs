package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(volumesCmd)
	rootCmd.AddCommand(snapshotsCmd)
	rootCmd.AddCommand(pvsCmd)
}

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "manage volumes",
	Long:  "manage volumes",
}

var snapshotsCmd = &cobra.Command{
	Use:   "snapshots",
	Short: "manage snapshots",
	Long:  "manage snapshots",
}

var pvsCmd = &cobra.Command{
	Use:   "pvs",
	Short: "manage pvs",
	Long:  "manage pvs",
}
