package cmd

import (
	"context"
	"fmt"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/superwomany/awsebs/internal/awsutils"
)

var vid string //aws://us-west-2a/vol-060c53922b1b43414

func init() {
	volumesCmd.AddCommand(encryptVCmd)
	encryptVCmd.Flags().StringVarP(&profile, "profile", "p", "sandbox-services", "aws profile to use")
	encryptVCmd.Flags().StringVarP(&vid, "vid", "v", "", "volume id return from pv, required")
}

var encryptVCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt volume from snapshot",
	Long:  "encrypt volume from snapshot",
	RunE:  encryptV,
	Example: `
go run main.go volumes encrypt -v aws://us-west-2a/vol-060c53922b1b43414
`,
}

func encryptV(cmd *cobra.Command, args []string) error {
	// aws://us-west-2a/vol-060c53922b1b43414
	filter := strings.Split(vid, "://")
	az := strings.Split(filter[1], "/")[0]
	id := strings.Split(filter[1], "/")[1]
	logger.WithFields(logger.Fields{
		"az":        az,
		"volume_id": id,
	}).Info("filter")
	vids := []string{id}

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
	if len(vid) < 1 {
		return fmt.Errorf("volume ids are required")
	}

	// get volume by id
	dl, err := awsutils.GetVolumeByID(context.Background(), spl, vids)
	if err != nil {
		return err
	}
	err = awsutils.PrintOutVList(dl)
	if err != nil {
		return err
	}
	// create snapshot from volume
	nspl, err := awsutils.CreateSnapshotsFromVolumeIds(context.Background(), spl, vids)
	if err != nil {
		return err
	}
	logger.WithFields(logger.Fields{
		"ids": nspl,
	}).Info("[encrypt_volume][CreateSnapshotsFromVolumeIds]")
	// check snapshot is ready
	newReady, err := awsutils.CompletedSnapshot(context.Background(), spl, nspl)
	if err != nil {
		return err
	}
	logger.WithFields(logger.Fields{
		"ready": newReady,
	}).Info("[encrypt_volume][CreateSnapshotsFromVolumeIds]")

	// copy snapshot
	if newReady {
		nids, err := awsutils.CopySnapshotByIDSWithEncryption(context.Background(), spl, nspl)
		if err != nil {
			return err
		}
		logger.WithFields(logger.Fields{
			"ids": nids,
		}).Info("[encrypt_volume][CopySnapshotByIDSWithEncryption]")
		// check snapshot is ready
		copyReady, err := awsutils.CompletedSnapshot(context.Background(), spl, nids)
		if err != nil {
			return err
		}
		logger.WithFields(logger.Fields{
			"ready": copyReady,
		}).Info("[encrypt_volume][CopySnapshotByIDSWithEncryption]")
		if copyReady {
			// use snapshot to create volume
			newVolumeID, err := awsutils.CreateVolumeFromSnapshot(context.Background(), spl, nids[0], az)
			if err != nil {
				return err
			}
			logger.WithFields(logger.Fields{
				"new_volume_id":      *newVolumeID,
				"original_volume_id": id,
				"old=>new":           id + "=>" + *newVolumeID,
			}).Info("[encrypt_volume][CreateVolumeFromSnapshot]")
		}

	}
	return nil
}
