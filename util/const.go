package util

import (
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

// DefaultPermission of the Kubernetes config file
const (
	DefaultPermission = 0664
)

// DefaultPollingDuration for CreateUpdate Call to AKS Cluster
const (
	DefaultPollingDuration = 1 * time.Hour
)

// DefaultKubeConfig is the location of Default Kubernetes config file
var (
	DefaultKubeConfig = clientcmd.RecommendedHomeFile
)
