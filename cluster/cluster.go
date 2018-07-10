// Package cluster provides operations on the Kubernetes cluster.
package cluster

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// New returns a new Service.
func New(clusterName string, s *session.Session) *Service {
	return &Service{
		ClusterName: clusterName,
		ec2:         ec2.New(s),
	}
}

// Service provides operations on the Kubernetes cluster.
type Service struct {
	ClusterName string
	DryRun      bool
	ec2         *ec2.EC2
}
