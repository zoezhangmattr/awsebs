package awsutils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	logger "github.com/sirupsen/logrus"
	"github.com/superwomany/awsebs/internal/utils"
)

func ListSnapshotByProfile(ctx context.Context, p ProfileData) ([]types.Snapshot, error) {
	logger.WithFields(logger.Fields{
		"profile": p,
	}).Info("[ListSnapshotByProfile]start.")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(p.Name))
	if err != nil {
		logger.Error("ListSnapshotByProfile]load config failed")
		return nil, err
	}
	smc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = p.Region
	})
	unEncryptedFile := types.Filter{
		Name:   aws.String("encrypted"),
		Values: []string{"false"},
	}
	el, err := smc.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		// MaxResults: 100,
		OwnerIds: []string{"self"},
		Filters:  []types.Filter{unEncryptedFile},
	})
	if err != nil {
		logger.WithFields(logger.Fields{
			"profile": p.Name,
		}).Errorf("[ListSnapshotByProfile] failed %v", err)
		return nil, err
	}
	logger.WithFields(logger.Fields{
		"profile": p.Name,
		"count":   len(el.Snapshots),
	}).Info("[ListSnapshotByProfile] done.")

	return el.Snapshots, nil
}

func ListUnEncryptedSnapshotByProfile(ctx context.Context, ec2Client *ec2.Client) ([]types.Snapshot, error) {
	logger.WithFields(logger.Fields{}).Info("[ListUnEncryptedSnapshotByProfile]start.")
	el, err := ec2Client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		MaxResults: aws.Int32(20),
		OwnerIds:   []string{"self"},
	})
	if err != nil {
		logger.WithFields(logger.Fields{}).Errorf("[ListUnEncryptedSnapshotByProfile] failed %v", err)
		return nil, err
	}
	logger.WithFields(logger.Fields{
		"count": len(el.Snapshots),
	}).Info("[ListUnEncryptedSnapshotByProfile] done.")

	return el.Snapshots, nil
}
func BatchCopySnapshotWithEncryption(ctx context.Context, pl ProfileData) error {
	logger.WithFields(logger.Fields{
		"profiles": pl,
	}).Info("[BatchCopySnapshotWithEncryption] start")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(pl.Name))
	if err != nil {
		logger.Error("BatchCopySnapshotWithEncryption]load config failed")
		return err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = pl.Region
	})
	dl, err := ListUnEncryptedSnapshotByProfile(context.Background(), ecc)
	if err != nil {
		return err
	}
	var pwg sync.WaitGroup
	pwg.Add(len(dl))
	errChan := make(chan error, len(dl))
	newSnapshots := []string{}
	for _, v := range dl {
		go func(ctx context.Context, ecc *ec2.Client, sp types.Snapshot) {
			defer pwg.Done()

			nid, err := CopySnapshotByParamWithEncryption(ctx, ecc, sp, pl.Region)
			if err != nil {
				nerr := fmt.Errorf("dl: %v, error: %v", sp, err)
				errChan <- nerr
				return
			}
			newSnapshots = append(newSnapshots, *nid)

		}(ctx, ecc, v)
	}
	pwg.Wait()
	close(errChan)
	if err := utils.ParseErrorsFromChannel(errChan); err != nil {
		logger.Errorf("we have errors %v, but it is ok to continue.", err)
		// return nil, err
	}
	logger.WithFields(logger.Fields{
		"profile": pl,
	}).Info("[BatchCopySnapshotWithEncryption] done")
	return nil

}

func CopySnapshotByParamWithEncryption(ctx context.Context, ec2Client *ec2.Client, params types.Snapshot, region string) (*string, error) {
	logger.WithFields(logger.Fields{}).Info("[CopySnapshotByParamWithEncryption]start.")
	tagSpec := types.TagSpecification{
		ResourceType: types.ResourceTypeSnapshot,
	}
	encryptedTag := types.Tag{
		Key:   aws.String("EncryptedByAutomation"),
		Value: aws.String("true"),
	}
	originalSpTag := types.Tag{
		Key:   aws.String("OriginalSnapshotID"),
		Value: aws.String(*params.SnapshotId),
	}
	if len(params.Tags) > 1 {
		isOrigial := false
		isEncrypted := false
		for _, i := range params.Tags {
			if *i.Key == "Name" {
				i.Value = aws.String("encrypted-" + *i.Value)
			}
			if *i.Key == "OriginalSnapshotID" {
				isOrigial = true
			}
			if *i.Key == "EncryptedByAutomation" {
				isEncrypted = true
			}
			tagSpec.Tags = append(tagSpec.Tags, i)
		}
		if !isOrigial {
			tagSpec.Tags = append(tagSpec.Tags, originalSpTag)
		}
		if !isEncrypted {
			tagSpec.Tags = append(tagSpec.Tags, encryptedTag)
		}

	} else {
		tagSpec.Tags = append(tagSpec.Tags, encryptedTag, originalSpTag)
	}
	res, err := ec2Client.CopySnapshot(ctx, &ec2.CopySnapshotInput{
		SourceSnapshotId:  params.SnapshotId,
		Encrypted:         aws.Bool(true),
		TagSpecifications: []types.TagSpecification{tagSpec},
		SourceRegion:      &region,
	})

	if err != nil {
		return nil, err
	}
	logger.WithFields(logger.Fields{"id": *res.SnapshotId}).Info("[CopySnapshotByParamWithEncryption]done")
	return res.SnapshotId, nil
}

func GetSnapshotByID(ctx context.Context, pl ProfileData, ids []string) ([]types.Snapshot, error) {
	logger.WithFields(logger.Fields{}).Info("[GetSnapshotByID]start.")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(pl.Name))
	if err != nil {
		logger.Error("[GetSnapshotByID]load config failed")
		return nil, err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = pl.Region
	})

	res, err := ecc.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		SnapshotIds: ids,
		OwnerIds:    []string{"self"},
	})

	if err != nil {
		return nil, err
	}
	for _, j := range res.Snapshots {
		logger.WithFields(logger.Fields{
			"encrypted":  *j.Encrypted,
			"volume":     *j.VolumeId,
			"snapshotid": *j.SnapshotId,
			// "tags":      j.Tags,
		}).Info("[GetSnapshotByID]done")
	}

	logger.WithFields(logger.Fields{
		"count": len(res.Snapshots),
	}).Info("[GetSnapshotByID]done")
	return res.Snapshots, nil
}

func CopySnapshotByIDSWithEncryption(ctx context.Context, pl ProfileData, ids []string) ([]string, error) {
	logger.WithFields(logger.Fields{}).Info("[CopySnapshotByIDSWithEncryption]start.")
	dl, err := GetSnapshotByID(ctx, pl, ids)
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(pl.Name))
	if err != nil {
		logger.Error("CopySnapshotByIDSWithEncryption]load config failed")
		return nil, err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = pl.Region
	})
	var pwg sync.WaitGroup
	pwg.Add(len(dl))
	errChan := make(chan error, len(dl))
	newSnapshots := []string{}
	for _, v := range dl {
		go func(ctx context.Context, ecc *ec2.Client, sp types.Snapshot) {
			defer pwg.Done()

			nid, err := CopySnapshotByParamWithEncryption(ctx, ecc, sp, pl.Region)
			if err != nil {
				nerr := fmt.Errorf("dl: %v, error: %v", sp, err)
				errChan <- nerr
				return
			}
			newSnapshots = append(newSnapshots, *nid)

		}(ctx, ecc, v)
	}
	pwg.Wait()
	close(errChan)
	if err := utils.ParseErrorsFromChannel(errChan); err != nil {
		logger.Errorf("we have errors %v, but it is ok to continue.", err)
		// return nil, err
	}
	logger.WithFields(logger.Fields{
		"profile": pl,
	}).Info("[CopySnapshotByIDSWithEncryption] done")

	return newSnapshots, nil
}

func CreateSnapshotsFromVolumeIds(ctx context.Context, p ProfileData, ids []string) ([]string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(p.Name))
	if err != nil {
		logger.Error("CreateSnapshotsFromVolume]load config failed")
		return nil, err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = p.Region
	})
	dl, err := GetVolumeByID(ctx, p, ids)
	if err != nil {
		return nil, err
	}

	var pwg sync.WaitGroup
	pwg.Add(len(dl))
	errChan := make(chan error, len(dl))
	newSnapshots := []string{}
	for _, v := range dl {
		go func(ctx context.Context, ecc *ec2.Client, vp types.Volume) {
			defer pwg.Done()
			tspec := types.TagSpecification{
				ResourceType: types.ResourceTypeSnapshot,
			}
			isaz := false
			if len(vp.Tags) > 0 {
				tspec.Tags = append(tspec.Tags, vp.Tags...)
				for _, t := range vp.Tags {
					if *t.Key == "volume-availability-zone" {
						isaz = true
					}
				}
			}
			attrTags := types.Tag{
				Key:   aws.String("volume-availability-zone"),
				Value: vp.AvailabilityZone,
			}
			if !isaz {
				tspec.Tags = append(tspec.Tags, attrTags)
			}

			res, err := ecc.CreateSnapshot(ctx, &ec2.CreateSnapshotInput{
				VolumeId:          vp.VolumeId,
				TagSpecifications: []types.TagSpecification{tspec},
			})
			if err != nil {
				nerr := fmt.Errorf("dl: %v, error: %v", vp, err)
				errChan <- nerr
				return
			}
			newSnapshots = append(newSnapshots, *res.SnapshotId)

			logger.WithFields(logger.Fields{
				"profile":    p,
				"snapshotid": *res.SnapshotId,
				"encrypted":  *res.Encrypted,
			}).Info("[CreateSnapshotsFromVolume] done")

		}(ctx, ecc, v)
	}
	pwg.Wait()
	close(errChan)
	if err := utils.ParseErrorsFromChannel(errChan); err != nil {
		logger.Errorf("we have errors %v, but it is ok to continue.", err)
		// return nil, err
	}
	logger.WithFields(logger.Fields{
		"profile": p,
	}).Info("[CreateSnapshotsFromVolume] done")

	return newSnapshots, nil
}

func CreateSnapshotFromVolume(ctx context.Context, p ProfileData, params types.Volume) (*string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(p.Name))
	if err != nil {
		logger.Error("CreateSnapshotFromVolume]load config failed")
		return nil, err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = p.Region
	})

	res, err := ecc.CreateSnapshot(ctx, &ec2.CreateSnapshotInput{
		VolumeId: aws.String(*params.VolumeId),
	})
	if err != nil {
		return nil, err
	}
	logger.WithFields(logger.Fields{
		"snapshotid": *res.SnapshotId,
		"volumeId":   *res.VolumeId,
		"encrypted":  res.Encrypted,
		"state":      res.State,
	}).Info("[CreateSnapshotFromVolume] done")
	return res.SnapshotId, nil
}

func CompletedSnapshot(ctx context.Context, pl ProfileData, ids []string) (bool, error) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return false, nil
		case <-ticker.C:
			b, err := SnapshotStatusCompletedByID(ctx, pl, ids)
			if err != nil {
				return false, err
			}
			if b {
				return b, nil
			}
		}
	}
}

func SnapshotStatusCompletedByID(ctx context.Context, pl ProfileData, ids []string) (bool, error) {
	logger.WithFields(logger.Fields{}).Info("[GetSnapshotStatusByID]start.")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(pl.Name))
	if err != nil {
		logger.Error("[GetSnapshotStatusByID]load config failed")
		return false, err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = pl.Region
	})

	res, err := ecc.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		SnapshotIds: ids,
		OwnerIds:    []string{"self"},
	})

	if err != nil {
		return false, err
	}
	count := 0

	for _, j := range res.Snapshots {
		logger.WithFields(logger.Fields{
			"encrypted":  *j.Encrypted,
			"volume":     *j.VolumeId,
			"snapshotid": *j.SnapshotId,
			"state":      j.State,
			// "tags":      j.Tags,
		}).Info("[GetSnapshotStatusByID]done")
		if j.State == types.SnapshotStateCompleted {
			count += 1
		}
	}
	if count == len(res.Snapshots) {
		return true, nil
	}

	return false, nil
}
