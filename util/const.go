package util

import (
	"k8s.io/client-go/tools/clientcmd"
)

// DefaultPermission of the Kubernetes config file
const (
	DefaultPermission = 0664
)

// DefaultKubeConfig is the location of Default Kubernetes config file
var (
	DefaultKubeConfig = clientcmd.RecommendedHomeFile
)
