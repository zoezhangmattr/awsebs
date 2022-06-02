package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

var ids []string

func init() {
	rootCmd.AddCommand(getSPCmd)
	getSPCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
	getSPCmd.Flags().StringArrayVarP(&ids, "ids", "i", nil, "snapshot ids")
}

var getSPCmd = &cobra.Command{
	Use:   "gets",
	Short: "get snapshots by id",
	Long:  "get snapshots by id",
	RunE:  getS,
	Example: `
go run main.go gets -i xxx -i xxx
	`,
}

func getS(cmd *cobra.Command, args []string) error {
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

	dl, err := awsutils.GetSnapshotByID(context.Background(), spl, ids)
	if err != nil {
		return err
	}
	err = awsutils.PrintOutSList(dl)
	if err != nil {
		return err
	}
	a, err := awsutils.CompletedSnapshot(context.Background(), spl, ids)
	if err != nil {
		return err
	}
	fmt.Print(a)
	return nil
}
