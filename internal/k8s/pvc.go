package k8s

import (
	"context"

	logger "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPVCs(ctx context.Context, c *kubernetes.Clientset) error {
	pvcl, err := c.CoreV1().PersistentVolumeClaims("all").List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.WithFields(logger.Fields{
		"pvc_count": len(pvcl.Items),
	}).Info("done")
	for _, j := range pvcl.Items {
		logger.WithFields(logger.Fields{
			"storageclass": j.Spec.StorageClassName,
			"name":         j.Name,
		}).Info("pvc detail")
	}
	return nil
}
