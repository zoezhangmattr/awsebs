package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

func init() {
	snapshotsCmd.AddCommand(copySPsCmd)
	snapshotsCmd.AddCommand(copySPCmd)
	copySPCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
	copySPCmd.Flags().StringArrayVarP(&ids, "ids", "i", nil, "snapshot ids")

	copySPsCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
}

var copySPsCmd = &cobra.Command{
	Use:   "copys",
	Short: "copys all unencypted snapshots with encrypted enabled",
	Long:  "copys all unencypted snapshots with encrypted enabled",
	RunE:  copySS,
}
var copySPCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy all unencypted snapshots with encrypted enabled by ids",
	Long:  "copy all unencypted snapshots with encrypted enabled by ids",
	RunE:  copyS,
}

func copySS(cmd *cobra.Command, args []string) error {
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

	err = awsutils.BatchCopySnapshotWithEncryption(context.Background(), spl)
	if err != nil {
		return err
	}

	return nil
}

func copyS(cmd *cobra.Command, args []string) error {
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
	if len(ids) == 0 {
		return fmt.Errorf("snapshot ids are required")

	}

	_, err = awsutils.CopySnapshotByIDSWithEncryption(context.Background(), spl, ids)
	if err != nil {
		return err
	}

	return nil
}
