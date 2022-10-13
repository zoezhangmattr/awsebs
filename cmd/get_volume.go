package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

func init() {
	volumesCmd.AddCommand(getVCmd)
	getVCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
	getVCmd.Flags().StringArrayVarP(&ids, "ids", "i", nil, "volume ids")
}

var getVCmd = &cobra.Command{
	Use:   "get",
	Short: "get volume by id",
	Long:  "get volume by id",
	RunE:  getV,
	Example: `
go run main.go volumes get -i xxx -i xxx
	`,
}

func getV(cmd *cobra.Command, args []string) error {
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

	dl, err := awsutils.GetVolumeByID(context.Background(), spl, ids)
	if err != nil {
		return err
	}
	err = awsutils.PrintOutVList(dl)
	if err != nil {
		return err
	}
	return nil
}
