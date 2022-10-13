package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

var profile string

func init() {
	snapshotsCmd.AddCommand(listSPCmd)
	listSPCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
}

var listSPCmd = &cobra.Command{
	Use:   "list",
	Short: "list all snapshots",
	Long:  "list all snapshots",
	RunE:  listS,
}

func listS(cmd *cobra.Command, args []string) error {
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
