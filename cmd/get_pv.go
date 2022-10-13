package cmd

import (
	"context"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/k8s"
)

var k8sc string

func init() {
	pvsCmd.AddCommand(getPVSCmd)
	getPVSCmd.Flags().StringVarP(&k8sc, "context", "c", "sandbox", "k8s context")
}

var getPVSCmd = &cobra.Command{
	Use:   "get",
	Short: "get all pv",
	Long:  "get all pv",
	RunE:  getPVS,
	Example: `
go run main.go pvs get
	`,
}

func getPVS(cmd *cobra.Command, args []string) error {

	ks := k8s.K8s{
		Context: k8sc,
	}
	c, err := ks.Auth()
	if err != nil {
		logger.Error("can't auth")
		return err
	}
	ctx := context.Background()
	err = k8s.GetPVs(ctx, c)
	if err != nil {
		return err
	}

	return nil
}
