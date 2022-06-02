package k8s

import (
	"context"

	logger "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPVs(ctx context.Context, c *kubernetes.Clientset) error {
	pvl, err := c.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.WithFields(logger.Fields{
		"pv_count": len(pvl.Items),
	}).Info("done")
	for _, j := range pvl.Items {
		logger.WithFields(logger.Fields{
			"volume_id":    j.Spec.AWSElasticBlockStore.VolumeID,
			"storageclass": j.Spec.StorageClassName,
			"claim":        j.Spec.ClaimRef.Name,
			"name":         j.Name,
		}).Info("pv detail")
	}
	return nil
}

// func CreatePV(ctx context.Context, c *kubernetes.Clientset, vid string) error {

// 	c.CoreV1().PersistentVolumes().Create(ctx, &v1.PersistentVolume{

// 	})
// 	return nil
// }
