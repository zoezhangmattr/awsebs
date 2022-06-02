package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

func init() {
	rootCmd.AddCommand(listVCmd)
	listVCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
}

var listVCmd = &cobra.Command{
	Use:   "lv",
	Short: "list all volumes",
	Long:  "list all volumes",
	RunE:  listV,
}

func listV(cmd *cobra.Command, args []string) error {
	pls, err := awsutils.LoadProfileConfig()
	if err != nil {
		return err
	}
	spl := awsutils.ProfileData{}
	vp := false
	for _, j := range pls {
		if j.Name == profile {
			spl = j
			vp = true
			break
		}
	}
	if !vp {
		return err
	}

	dl, err := awsutils.ListSnapshotByProfile(context.Background(), spl)
	if err != nil {
		return err
	}
	err = awsutils.PrintOutSList(dl)
	if err != nil {
		return err
	}
	return nil

}
