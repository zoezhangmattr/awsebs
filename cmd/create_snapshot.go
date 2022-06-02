package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

func init() {
	rootCmd.AddCommand(createSPCmd)
	createSPCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
	createSPCmd.Flags().StringArrayVarP(&ids, "ids", "i", nil, "volume ids")

}

var createSPCmd = &cobra.Command{
	Use:   "creates",
	Short: "create snapshot from volume id",
	Long:  "create snapshot from volume id",
	RunE:  createS,
	Example: `
go run main.go creates -i vol-xxxxxx
	`,
}

func createS(cmd *cobra.Command, args []string) error {
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
	if len(ids) < 1 {
		return fmt.Errorf("volume ids are required")
	}

	_, err = awsutils.CreateSnapshotsFromVolumeIds(context.Background(), spl, ids)
	if err != nil {
		return err
	}
	return nil
}
