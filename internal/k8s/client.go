package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
)

// K8sProvider defines a K8s clientset Auth functionality
type K8sProvider interface {
	Auth() (*kubernetes.Clientset, error)
}

// K8s representation for a k8s client set
type K8s struct {
	Context string
}

// Auth fetches the k8s client based on incluster run vs local run
func (kp K8s) Auth() (*kubernetes.Clientset, error) {

	c, err := kp.GetClient()
	if err != nil {
		return nil, errors.Wrap(err, "Could not Authenticate to K8s")
	}

	return c, nil
}

// GetClient returns a K8s clientset
func (kp K8s) GetClient() (*kubernetes.Clientset, error) {

	var kc *string
	if home := kp.homeDir(); home != "" {
		fp := filepath.Join(home, ".kube", "config")
		kc = &fp
	} else {
		return nil, fmt.Errorf("kubeconfig is not found")
	}

	logger.WithFields(logger.Fields{
		"kubeconfig": *kc,
		"context":    kp.Context,
	}).Info("k8s get client")

	var c *rest.Config
	var err error

	if kp.Context != "" {
		configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kc}
		configOverrides := &clientcmd.ConfigOverrides{CurrentContext: kp.Context}

		c, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
		if err != nil {
			return nil, err
		}
	} else {
		c, err = clientcmd.BuildConfigFromFlags("", *kc)
		if err != nil {
			return nil, errors.Wrap(err, "Could not build K8s config from flag")
		}
	}

	logger.WithFields(logger.Fields{
		"kubeconfig":    *kc,
		"host":          c.Host,
		"ExecProvider:": c.ExecProvider,
	}).Info("k8s client config")

	cs, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to create kubernetes client")
	}

	return cs, nil
}

func (kp K8s) homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
