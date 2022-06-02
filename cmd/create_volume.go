package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

var az string
var sid string

func init() {
	rootCmd.AddCommand(createVCmd)
	createVCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
	createVCmd.Flags().StringVarP(&sid, "sid", "i", "", "snapshot id, required")
	createVCmd.Flags().StringVarP(&az, "az", "a", "", "volume availability zone, required")
}

var createVCmd = &cobra.Command{
	Use:   "createv",
	Short: "create volume from snapshot",
	Long:  "create volume from snapshot",
	RunE:  createV,
	Example: `
go run main.go createv -i snap-xxxxxx -a ap-southeast-2b
`,
}

func createV(cmd *cobra.Command, args []string) error {
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
	if len(sid) < 1 {
		return fmt.Errorf("volume ids are required")
	}
	if len(az) < 1 {
		return fmt.Errorf("az is required")
	}

	_, err = awsutils.CreateVolumeFromSnapshot(context.Background(), spl, sid, az)
	if err != nil {
		return err
	}

	return nil
}
