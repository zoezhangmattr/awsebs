package awsutils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	logger "github.com/sirupsen/logrus"
)

func GetVolumeByID(ctx context.Context, pl ProfileData, ids []string) ([]types.Volume, error) {
	logger.WithFields(logger.Fields{}).Info("[GetVolumeByID]start.")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(pl.Name))
	if err != nil {
		logger.Error("[GetVolumeByID]load config failed")
		return nil, err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = pl.Region
	})

	res, err := ecc.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
		VolumeIds: ids,
	})

	if err != nil {
		return nil, err
	}
	for _, j := range res.Volumes {
		logger.WithFields(logger.Fields{
			"encrypted": *j.Encrypted,
			"volume":    *j.VolumeId,
			"az":        *j.AvailabilityZone,
			"state":     j.State,
		}).Info("[GetVolumeByID]done")
		for _, t := range j.Tags {
			logger.WithFields(logger.Fields{
				"key":   *t.Key,
				"value": *t.Value,
			}).Info("[GetVolumeByID] tags")
		}
	}

	logger.WithFields(logger.Fields{
		"count": len(res.Volumes),
	}).Info("[GetVolumeByID]done")

	return res.Volumes, nil
}

func CreateVolumeFromSnapshot(ctx context.Context, p ProfileData, id string, az string) (*string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(p.Name))
	if err != nil {
		logger.Error("CreateVolumeFromSnapshot]load config failed")
		return nil, err
	}
	ecc := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = p.Region
	})
	si, err := GetSnapshotByID(ctx, p, []string{id})
	if err != nil {
		return nil, err
	}
	if len(si) < 1 {
		return nil, err
	}
	sts := si[0].Tags
	customTag := types.Tag{
		Key:   aws.String("Automation"),
		Value: aws.String("true"),
	}
	isct := false
	for _, t := range sts {
		if *t.Key == "Automation" {
			isct = true
			break
		}
	}
	if !isct {
		sts = append(sts, customTag)
	}
	tspec := types.TagSpecification{
		ResourceType: types.ResourceTypeVolume,
		Tags:         sts,
	}
	res, err := ecc.CreateVolume(ctx, &ec2.CreateVolumeInput{
		AvailabilityZone:  aws.String(az),
		SnapshotId:        aws.String(id),
		TagSpecifications: []types.TagSpecification{tspec},
		VolumeType:        types.VolumeTypeGp3, // make it gp3
	})
	if err != nil {
		return nil, err
	}
	logger.WithFields(logger.Fields{
		"snapshotid":       id,
		"availabilityzone": *res.AvailabilityZone,
		"volumeId":         *res.VolumeId,
		"encrypted":        *res.Encrypted,
		"createdTime":      res.CreateTime,
		"pv-volumeid":      "aws://" + *res.AvailabilityZone + "/" + *res.VolumeId,
	}).Info("[CreateVolumeFromSnapshot] done")
	return res.VolumeId, nil
}
